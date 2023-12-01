package match

import (
	"time"

	"github.com/google/uuid"
)

type EnqueUserRes struct {
	ID        uuid.UUID `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
