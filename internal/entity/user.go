package entity

import "time"

type User struct {
	ID          string `gorm:"primaryKey"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	PhoneNumber string `gorm:"unique;not null"`
	Address     string `gorm:"not null"`
	Pin         string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Balance     UserBalance `gorm:"foreignKey:UserID"`
}

type UserBalance struct {
	UserID  string `gorm:"primaryKey"`
	Balance float64
}
