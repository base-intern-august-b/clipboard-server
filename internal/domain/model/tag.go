package model

import (
	"time"
)

type Tag struct {
	HashedTag string    `db:"hashed_tag" json:"hashed_tag"`
	RawTag    string    `db:"raw_tag" json:"raw_tag"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
