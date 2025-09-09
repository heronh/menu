package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name               string `gorm:"not null"`
	Email              string `gorm:"unique;not null"`
	Password           string `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CompanyID          *uint      // Pointer to allow null
	Company            Company    `gorm:"foreignKey:CompanyID"`
	AuthoredCategories []Category `gorm:"foreignKey:AuthorID"`
	AuthoredDishes     []Dish     `gorm:"foreignKey:AuthorLastChangeID"`
	InsertedImages     []Image    `gorm:"foreignKey:InsertedByID"`
	SentMessages       []Message  `gorm:"foreignKey:SenderID"`
	ReceivedMessages   []Message  `gorm:"foreignKey:RecipientID"`
	Todo               []Todo     `gorm:"foreignKey:UserID"`
	PrivilegeID        uint       `gorm:"not null"`
	Privilege          Privilege  `gorm:"foreignKey:PrivilegeID"`
}
