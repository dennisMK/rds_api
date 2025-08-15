package worker

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Job represents a unit of work
type Job struct {
	ID       string
	Type     string
	Payload  interface{}
	Retries  int
	MaxRetries int
	CreatedAt time.Time
}

// JobResult represents the result of a job execution
type JobResult struct {
	JobID     string
	Success   bool
	Error     error
	Duration  time.Duration
	CompletedAt time.Time
}

// JobHandler defines the interface for job handlers
type JobHandler interface {
	Handle(ctx context.Context, job *Job) error
	GetJobType() string
}

// WorkerPool manages a pool of workers for concurrent job processing
type WorkerPool struct {
	workers     int
	jobQueue    chan *Job
	resultQueue chan *JobResult
	quit        chan bool
	wg          sync.WaitGroup
	handlers    map[string]JobHandler
	logger      *logrus.Logger
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, queueSize int, logger *logrus.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workers:     workers,
		jobQueue:    make(chan *Job, queueSize),
		resultQueue: make(chan *JobResult, queueSize),
		quit:        make(chan bool),
		handlers:    make(map[string]JobHandler),
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// RegisterHandler registers a job handler for a specific job type
func (wp *WorkerPool) RegisterHandler(handler JobHandler) {
	wp.handlers[handler.GetJobType()] = handler
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	wp.logger.Infof("Starting worker pool with %d workers", wp.workers)
	
	// Start workers
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	
	// Start result processor
	go wp.processResults()
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.logger.Info("Stopping worker pool...")
	
	close(wp.quit)
	wp.cancel()
	wp.wg.Wait()
	
	close(wp.jobQueue)
	close(wp.resultQueue)
	
	wp.logger.Info("Worker pool stopped")
}

// SubmitJob submits a job to the worker pool
func (wp *WorkerPool) SubmitJob(job *Job) error {
	select {
	case wp.jobQueue <- job:
		wp.logger.WithFields(logrus.Fields{
			"job_id":   job.ID,
			"job_type": job.Type,
		}).Debug("Job submitted to queue")
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return ErrQueueFull
	}
}

// worker processes jobs from the job queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	wp.logger.WithField("worker_id", id).Debug("Worker started")
	
	for {
		select {
		case job := <-wp.jobQueue:
			if job == nil {
				return
			}
			wp.processJob(id, job)
			
		case <-wp.quit:
			wp.logger.WithField("worker_id", id).Debug("Worker stopping")
			return
		}
	}
}

// processJob processes a single job
func (wp *WorkerPool) processJob(workerID int, job *Job) {
	start := time.Now()
	
	logger := wp.logger.WithFields(logrus.Fields{
		"worker_id": workerID,
		"job_id":    job.ID,
		"job_type":  job.Type,
	})
	
	logger.Debug("Processing job")
	
	// Get handler for job type
	handler, exists := wp.handlers[job.Type]
	if !exists {
		logger.Error("No handler found for job type")
		wp.resultQueue <- &JobResult{
			JobID:       job.ID,
			Success:     false,
			Error:       ErrNoHandler,
			Duration:    time.Since(start),
			CompletedAt: time.Now(),
		}
		return
	}
	
	// Execute job with timeout
	ctx, cancel := context.WithTimeout(wp.ctx, 30*time.Second)
	defer cancel()
	
	err := handler.Handle(ctx, job)
	duration := time.Since(start)
	
	result := &JobResult{
		JobID:       job.ID,
		Success:     err == nil,
		Error:       err,
		Duration:    duration,
		CompletedAt: time.Now(),
	}
	
	if err != nil {
		logger.WithError(err).Error("Job failed")
		
		// Retry logic
		if job.Retries < job.MaxRetries {
			job.Retries++
			logger.WithField("retry_count", job.Retries).Info("Retrying job")
			
			// Exponential backoff
			backoff := time.Duration(job.Retries*job.Retries) * time.Second
			time.AfterFunc(backoff, func() {
				wp.SubmitJob(job)
			})
			return
		}
		
		logger.Error("Job failed after max retries")
	} else {
		logger.WithField("duration", duration).Debug("Job completed successfully")
	}
	
	// Send result
	select {
	case wp.resultQueue <- result:
	default:
		logger.Warn("Result queue full, dropping result")
	}
}

// processResults processes job results
func (wp *WorkerPool) processResults() {
	for result := range wp.resultQueue {
		wp.logger.WithFields(logrus.Fields{
			"job_id":   result.JobID,
			"success":  result.Success,
			"duration": result.Duration,
		}).Info("Job result processed")
		
		// Here you could store results in database, send notifications, etc.
	}
}

// GetStats returns worker pool statistics
func (wp *WorkerPool) GetStats() WorkerPoolStats {
	return WorkerPoolStats{
		Workers:        wp.workers,
		QueuedJobs:     len(wp.jobQueue),
		QueueCapacity:  cap(wp.jobQueue),
		PendingResults: len(wp.resultQueue),
	}
}

// WorkerPoolStats represents worker pool statistics
type WorkerPoolStats struct {
	Workers        int `json:"workers"`
	QueuedJobs     int `json:"queued_jobs"`
	QueueCapacity  int `json:"queue_capacity"`
	PendingResults int `json:"pending_results"`
}

// Custom errors
var (
	ErrQueueFull  = fmt.Errorf("job queue is full")
	ErrNoHandler  = fmt.Errorf("no handler found for job type")
)
