package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:created_at_index;not_null"`
	Tag       string    `gorm:"index:tag_index"`
	Text      string
}

type Messages []Message
