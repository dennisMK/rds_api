package concurrent

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Stage represents a processing stage in the pipeline
type Stage[T any] func(ctx context.Context, input <-chan T, output chan<- T) error

// Pipeline represents a concurrent processing pipeline
type Pipeline[T any] struct {
	stages []Stage[T]
	logger *logrus.Logger
}

// NewPipeline creates a new processing pipeline
func NewPipeline[T any](logger *logrus.Logger) *Pipeline[T] {
	return &Pipeline[T]{
		logger: logger,
	}
}

// AddStage adds a processing stage to the pipeline
func (p *Pipeline[T]) AddStage(stage Stage[T]) *Pipeline[T] {
	p.stages = append(p.stages, stage)
	return p
}

// Process processes items through the pipeline
func (p *Pipeline[T]) Process(ctx context.Context, items []T, bufferSize int) ([]T, error) {
	if len(p.stages) == 0 {
		return items, nil
	}

	p.logger.WithFields(logrus.Fields{
		"items":  len(items),
		"stages": len(p.stages),
	}).Info("Starting pipeline processing")

	// Create channels for each stage
	channels := make([]chan T, len(p.stages)+1)
	for i := range channels {
		channels[i] = make(chan T, bufferSize)
	}

	// Start all stages
	var wg sync.WaitGroup
	errChan := make(chan error, len(p.stages))

	for i, stage := range p.stages {
		wg.Add(1)
		go func(stageIndex int, stageFunc Stage[T]) {
			defer wg.Done()
			defer close(channels[stageIndex+1])

			p.logger.WithField("stage", stageIndex).Debug("Starting pipeline stage")

			if err := stageFunc(ctx, channels[stageIndex], channels[stageIndex+1]); err != nil {
				p.logger.WithError(err).WithField("stage", stageIndex).Error("Pipeline stage failed")
				errChan <- err
			}

			p.logger.WithField("stage", stageIndex).Debug("Pipeline stage completed")
		}(i, stage)
	}

	// Feed input items
	go func() {
		defer close(channels[0])
		for _, item := range items {
			select {
			case channels[0] <- item:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Collect output
	var results []T
	go func() {
		for item := range channels[len(channels)-1] {
			results = append(results, item)
		}
	}()

	// Wait for completion
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	p.logger.WithField("results", len(results)).Info("Pipeline processing completed")
	return results, nil
}
