package kafka

import (
	"backend/internal/chat"
	"log"
	"sync"
	"time"
)

type BatchProcessor struct {
	mu        sync.Mutex
	buffer    []chat.Message
	maxSize   int
	maxWait   time.Duration
	saveFunc  func([]chat.Message) error
	waitGroup sync.WaitGroup
}

func NewBatchProcessor(maxSize int, maxWait time.Duration, saveFunc func([]chat.Message) error) *BatchProcessor {
	bp := &BatchProcessor{
		buffer:   make([]chat.Message, 0, maxSize),
		maxSize:  maxSize,
		maxWait:  maxWait,
		saveFunc: saveFunc,
	}
	bp.waitGroup.Add(1)
	go bp.run()
	return bp
}

func (bp *BatchProcessor) AddMessage(msg chat.Message) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.buffer = append(bp.buffer, msg)
	if len(bp.buffer) >= bp.maxSize {
		bp.flush()
	}
}

func (bp *BatchProcessor) run() {
	defer bp.waitGroup.Done()
	ticker := time.NewTicker(bp.maxWait)
	defer ticker.Stop()

	for range ticker.C {
		bp.mu.Lock()
		bp.flush()
		bp.mu.Unlock()
	}
}

func (bp *BatchProcessor) flush() {
	if len(bp.buffer) > 0 {
		if err := bp.saveFunc(bp.buffer); err != nil {
			// Handle error...
			log.Printf("Failed to flush messages: %v", err)
		}
		bp.buffer = bp.buffer[:0] // Clear buffer
	}
}

func (bp *BatchProcessor) Stop() {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.flush()
	bp.waitGroup.Wait()
}
