package kafka

import (
	"backend/internal/chat"
	"database/sql"
	"log"
	"time"
)

const (
	batchSize     = 500
	batchInterval = 2 * time.Second
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

// Batch insert chat messages into Postgres using a transaction
func NewInsertFunc(db *sql.DB) func([]chat.Message) error {
	return func(messages []chat.Message) error {
		// Begin a transaction
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		// Prepare the insert statement
		stmt, err := tx.Prepare("INSERT INTO messages (chat_room_id, sender_id, content, media_url, created_at) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			tx.Rollback() // Don't forget to rollback in case of an error
			return err
		}
		defer stmt.Close()

		// Iterate over the messages and execute the insert statement
		for _, msg := range messages {
			_, err := stmt.Exec(msg.ChatRoomID, msg.SenderID, msg.Content, msg.MediaURL, msg.CreatedAt)
			if err != nil {
				tx.Rollback() // Rollback the transaction in case of an error
				return err
			}
		}

		// Commit the transaction
		return tx.Commit()
	}
}
