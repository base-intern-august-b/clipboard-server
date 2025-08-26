package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// User はユーザーを表すドメインモデル
type User struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	UserName  string    `db:"user_name" json:"user_name"`
	Nickname  string    `db:"nickname" json:"nickname"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type RequestCreateUser struct {
	UserName string `db:"user_name" json:"user_name"`
	Password string `db:"password" json:"password"`
	Nickname string `db:"nickname" json:"nickname"`
	Status   string `db:"status" json:"status"`
}

type RequestGetUserBatch struct {
	UserNames []string `json:"user_names"`
}

type RequestPatchUser struct {
	UserName *string `json:"user_name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Status   *string `json:"status,omitempty"`
}

type RequestChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
