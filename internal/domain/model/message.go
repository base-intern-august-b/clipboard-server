package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	MessageID uuid.UUID `db:"message_id" json:"message_id"`
	ChannelID uuid.UUID `db:"channel_id" json:"channel_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type RequestCreateMessage struct {
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
}

type RequestPatchMessage struct {
	Content *string `json:"content,omitempty"`
}

type RequestGetMessages struct {
	ChannelID uuid.UUID `json:"channel_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

type RequestGetMessagesInDuration struct {
	ChannelID uuid.UUID `json:"channel_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
