package concurrent

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// BatchProcessor processes items in batches with concurrency control
type BatchProcessor[T any] struct {
	batchSize   int
	maxWorkers  int
	timeout     time.Duration
	processor   func(ctx context.Context, batch []T) error
	logger      *logrus.Logger
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor[T any](
	batchSize int,
	maxWorkers int,
	timeout time.Duration,
	processor func(ctx context.Context, batch []T) error,
	logger *logrus.Logger,
) *BatchProcessor[T] {
	return &BatchProcessor[T]{
		batchSize:  batchSize,
		maxWorkers: maxWorkers,
		timeout:    timeout,
		processor:  processor,
		logger:     logger,
	}
}

// Process processes items in batches concurrently
func (bp *BatchProcessor[T]) Process(ctx context.Context, items []T) error {
	if len(items) == 0 {
		return nil
	}

	// Create batches
	batches := bp.createBatches(items)
	
	// Create worker pool
	semaphore := make(chan struct{}, bp.maxWorkers)
	var wg sync.WaitGroup
	errChan := make(chan error, len(batches))

	bp.logger.WithFields(logrus.Fields{
		"total_items": len(items),
		"batches":     len(batches),
		"batch_size":  bp.batchSize,
		"max_workers": bp.maxWorkers,
	}).Info("Starting batch processing")

	// Process batches concurrently
	for i, batch := range batches {
		wg.Add(1)
		go func(batchIndex int, batchItems []T) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Create context with timeout
			batchCtx, cancel := context.WithTimeout(ctx, bp.timeout)
			defer cancel()

			bp.logger.WithFields(logrus.Fields{
				"batch_index": batchIndex,
				"batch_size":  len(batchItems),
			}).Debug("Processing batch")

			start := time.Now()
			if err := bp.processor(batchCtx, batchItems); err != nil {
				bp.logger.WithError(err).WithField("batch_index", batchIndex).Error("Batch processing failed")
				errChan <- err
				return
			}

			bp.logger.WithFields(logrus.Fields{
				"batch_index": batchIndex,
				"duration":    time.Since(start),
			}).Debug("Batch processed successfully")
		}(i, batch)
	}

	// Wait for all batches to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		bp.logger.WithField("error_count", len(errors)).Error("Batch processing completed with errors")
		return errors[0] // Return first error
	}

	bp.logger.WithField("total_items", len(items)).Info("Batch processing completed successfully")
	return nil
}

// createBatches splits items into batches
func (bp *BatchProcessor[T]) createBatches(items []T) [][]T {
	var batches [][]T
	
	for i := 0; i < len(items); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	
	return batches
}
