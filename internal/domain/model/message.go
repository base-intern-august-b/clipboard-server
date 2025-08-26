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

type MessageDetail struct {
	Message
	User *User `json:"user"`
	Tags []*Tag `json:"tags"`
}

type RequestCreateMessage struct {
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
}

type RequestPatchMessage struct {
	Content *string `json:"content,omitempty"`
}
