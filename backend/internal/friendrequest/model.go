package friendrequest

import (
	"time"

	"github.com/google/uuid"
)

type FriendRequestStatus string

const (
	Pending  FriendRequestStatus = "PENDING"
	Accepted FriendRequestStatus = "ACCEPTED"
	Declined FriendRequestStatus = "DECLINED"
)

type FriendRequest struct {
	ID          uuid.UUID           `json:"id"`
	SenderID    uuid.UUID           `json:"sender_id"`
	RecipientID uuid.UUID           `json:"recipient_id"`
	Status      FriendRequestStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
}
