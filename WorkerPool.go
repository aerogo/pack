package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// WorkerPoolResults maps jobs to their corresponding results.
type WorkerPoolResults map[interface{}]interface{}

// WorkerPool ...
type WorkerPool struct {
	jobCount    uint64
	jobs        chan interface{}
	done        chan bool
	results     WorkerPoolResults
	resultsLock sync.RWMutex
}

// NewWorkerPool ...
func NewWorkerPool(work func(interface{}) interface{}) *WorkerPool {
	pool := new(WorkerPool)
	pool.jobs = make(chan interface{}, 4096)
	pool.done = make(chan bool, 4096)
	pool.results = make(WorkerPoolResults)

	for w := 1; w <= runtime.NumCPU(); w++ {
		go func() {
			for job := range pool.jobs {
				result := work(job)

				pool.resultsLock.Lock()
				pool.results[job] = result
				pool.resultsLock.Unlock()

				pool.done <- true
			}
		}()
	}

	return pool
}

// Queue ...
func (pool *WorkerPool) Queue(job interface{}) {
	pool.jobs <- job
	atomic.AddUint64(&pool.jobCount, 1)
}

// Wait ...
func (pool *WorkerPool) Wait() WorkerPoolResults {
	jobCount := atomic.LoadUint64(&pool.jobCount)

	for i := uint64(0); i < jobCount; i++ {
		<-pool.done
	}

	return pool.results
}
