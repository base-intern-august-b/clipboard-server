package model

import "time"

type Channel struct {
	ChannelID   int64     `db:"channel_id" json:"channel_id"`
	ChannelName string    `db:"channel_name" json:"channel_name"`
	DisplayName string    `db:"display_name" json:"display_name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type RequestCreateChannel struct {
	ChannelName string `db:"channel_name" json:"channel_name"`
	DisplayName string `db:"display_name" json:"display_name"`
	Description string `db:"description" json:"description"`
}

type RequestPatchChannel struct {
	ChannelName *string `json:"channel_name,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
}
