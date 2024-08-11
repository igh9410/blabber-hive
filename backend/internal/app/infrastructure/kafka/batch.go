package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/igh9410/blabber-hive/backend/internal/chat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BatchProcessor struct {
	InsertFunc func([]chat.Message) error
	BatchSize  int
	BatchTime  time.Duration
	Messages   []chat.Message
	Ticker     *time.Ticker
	TickerStop chan bool
	Lock       sync.Mutex
}

func NewBatchProcessor(insertFunc func([]chat.Message) error, batchSize int, batchTime time.Duration) *BatchProcessor {
	bp := &BatchProcessor{
		InsertFunc: insertFunc,
		BatchSize:  batchSize,
		BatchTime:  batchTime,
		Ticker:     time.NewTicker(batchTime),
		TickerStop: make(chan bool),
	}
	go bp.run()
	return bp
}

func (bp *BatchProcessor) run() {
	for {
		select {
		case <-bp.Ticker.C:
			bp.flush()
		case <-bp.TickerStop:
			return
		}
	}
}

func (bp *BatchProcessor) AddMessage(message chat.Message) {
	bp.Messages = append(bp.Messages, message)
	if len(bp.Messages) >= bp.BatchSize {
		bp.flush()
	}
}

func (bp *BatchProcessor) flush() {
	bp.Lock.Lock()
	messagesToInsert := bp.Messages
	bp.Messages = nil
	bp.Lock.Unlock()

	if len(messagesToInsert) > 0 {
		err := bp.InsertFunc(messagesToInsert)
		if err != nil {
			log.Printf("Failed to insert messages into Postgres: %s", err)
			return
			// Handle error
		}
	}
}

func (bp *BatchProcessor) Stop() {
	bp.Ticker.Stop()
	bp.TickerStop <- true
	bp.flush() // Flush any remaining messages
}

// / Bulk insert chat messages into Postgres using a transaction
func NewInsertFunc(db *pgxpool.Pool) func([]chat.Message) error {
	return func(messages []chat.Message) error {
		// Use context with a timeout for the transaction
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Begin a transaction
		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Prepare the insert statement with bulk support
		valueStrings := make([]string, 0, len(messages))
		valueArgs := make([]interface{}, 0, len(messages)*5)
		for i, msg := range messages {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
			valueArgs = append(valueArgs, msg.ChatRoomID, msg.SenderID, msg.Content, msg.MediaURL, msg.CreatedAt)
		}

		stmt := fmt.Sprintf("INSERT INTO messages (chat_room_id, sender_id, content, media_url, created_at) VALUES %s",
			strings.Join(valueStrings, ","))

		// Execute the bulk insert statement within the transaction
		_, err = tx.Exec(ctx, stmt, valueArgs...)
		if err != nil {
			// Rollback the transaction in case of an error
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %s", rollbackErr)
				return rollbackErr
			}
			return fmt.Errorf("failed to execute bulk insert: %w", err)
		}

		// Commit the transaction
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}
