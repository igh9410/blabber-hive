package kafka

import (
	"backend/internal/chat"
	"log"
	"time"
)

const (
	batchSize     = 100             // Change as per your requirement
	batchInterval = 5 * time.Second // Change as per your requirement
)

type BatchProcessor struct {
	Messages    []chat.Message
	InsertFunc  func([]chat.Message) error // Function to insert messages into Postgres
	BatchSize   int
	BatchTicker *time.Ticker
}

func NewBatchProcessor(insertFunc func([]chat.Message) error, batchSize int, batchInterval time.Duration) *BatchProcessor {
	bp := &BatchProcessor{
		Messages:    make([]chat.Message, 0, batchSize),
		InsertFunc:  insertFunc,
		BatchSize:   batchSize,
		BatchTicker: time.NewTicker(batchInterval),
	}

	go bp.run()
	return bp
}

func (bp *BatchProcessor) AddMessage(message chat.Message) {
	bp.Messages = append(bp.Messages, message)
	if len(bp.Messages) >= bp.BatchSize {
		bp.flush()
	}
}

func (bp *BatchProcessor) run() {
	for range bp.BatchTicker.C {
		bp.flush()
	}
}

func (bp *BatchProcessor) flush() {
	if len(bp.Messages) == 0 {
		return
	}

	// Insert messages into PostgreSQL
	if err := bp.InsertFunc(bp.Messages); err != nil {
		log.Printf("Failed to insert messages: %v", err)
		return
	}

	// Clear messages
	bp.Messages = bp.Messages[:0]
}
