package models

import (
	"time"

	"gorm.io/gorm"
)

// Company represents a company
type Company struct {
	gorm.Model
	Name         string `gorm:"unique;not null"`
	CEP          string
	Street       string
	Number       string
	Neighborhood string
	City         string
	State        string
	Active       bool   `gorm:"default:true"`
	CNPJ         string `gorm:"unique;not null"`
	Level        int    `gorm:"type:integer;check:level >= 0 AND level <= 100"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Users        []User    `gorm:"foreignKey:CompanyID"`
	Dishes       []Dish    `gorm:"foreignKey:CompanyID"`
	Images       []Image   `gorm:"foreignKey:CompanyID"`
	Messages     []Message `gorm:"foreignKey:SenderCompanyID"`
}

// Message represents a message between users
type Message struct {
	gorm.Model
	Text            string `gorm:"not null"`
	Subject         string
	SenderID        uint    `gorm:"not null"`
	Sender          User    `gorm:"foreignKey:SenderID"`
	SenderCompanyID uint    // Can be null if sender is SU or not associated with a company for the message
	SenderCompany   Company `gorm:"foreignKey:SenderCompanyID"`
	RecipientID     uint    `gorm:"not null"`
	Recipient       User    `gorm:"foreignKey:RecipientID"`
	SendDate        time.Time
}
