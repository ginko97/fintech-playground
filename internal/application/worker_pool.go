package application

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ginko97/fintech-playground/internal/domain"
)

type WorkerPool struct {
	workers      int
	jobQueue     chan *domain.Transaction
	stateMachine *TransactionStateMachine
	wg           sync.WaitGroup
}

func NewWorkerPool(workers int, sm *TransactionStateMachine) *WorkerPool {
	return &WorkerPool{
		workers:      workers,
		jobQueue:     make(chan *domain.Transaction, 1000),
		stateMachine: sm,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	log.Printf("Worker pool started with %d workers", wp.workers)
}

// Submit a transaction for background processing
func (wp *WorkerPool) Submit(tx *domain.Transaction) {
	wp.jobQueue <- tx
}

// worker processes jobs
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for tx := range wp.jobQueue {
		log.Printf("[Worker %d] Processing tx %s (status: %s)", id, tx.ID, tx.Status)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		// Simulate external processing (PSP call, etc.)
		err := wp.processTransaction(ctx, tx)
		cancel()

		if err != nil {
			log.Printf("[Worker %d] Failed tx %s: %v", id, tx.ID, err)
			// TODO: push to DLQ later
		}
	}
}

func (wp *WorkerPool) processTransaction(ctx context.Context, tx *domain.Transaction) error {
	// Example: move state forward
	if tx.Status == domain.StatusPending {
		_ = wp.stateMachine.UpdateState(ctx, tx.ID, domain.StatusProcessing, "started_processing")
		// Simulate work
		time.Sleep(500 * time.Millisecond)
		_ = wp.stateMachine.UpdateState(ctx, tx.ID, domain.StatusCompleted, "success")
	}
	return nil
}

// Shutdown gracefully
func (wp *WorkerPool) Shutdown() {
	close(wp.jobQueue)
	wp.wg.Wait()
	log.Println("Worker pool shutdown completed")
}
