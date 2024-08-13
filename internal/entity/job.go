package entity

import (
	"time"
)

// Job represents a background job for processing tasks
type Job struct {
	ID          uint      `gorm:"primaryKey"`
	Message     string    `gorm:"not null"`
	Status      string    `gorm:"not null"`
	QueueAt     time.Time `gorm:"autoCreateTime"`
	CompletedAt time.Time `gorm:""`
}

const (
	StatusQueued    = "QUEUED"
	StatusCompleted = "COMPLETED"
)
