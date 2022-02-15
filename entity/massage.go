package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index:created_at_index;not_null" json:"created_at"`
	Tag       string    `gorm:"index:tag_index" json:"tag"`
	Text      string    `json:"text"`
}

type Messages []*Message
