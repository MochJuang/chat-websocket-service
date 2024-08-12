package entity

import "time"

type Transaction struct {
	ID              string  `gorm:"primaryKey"`
	UserID          string  `gorm:"not null"`
	Amount          float64 `gorm:"not null"`
	TransactionType string  `gorm:"not null"`
	Remarks         string
	CreatedAt       time.Time
}

const (
	DEBIT  = "DEBIT"
	CREDIT = "CREDIT"
)
