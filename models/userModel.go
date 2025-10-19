package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name             string    `gorm:"not null"`
	Email            string    `gorm:"unique;not null"`
	Password         string    `gorm:"not null"`
	CompanyID        *uint     // Pointer to allow null
	Company          Company   `gorm:"foreignKey:CompanyID"`
	Sections         []Section `gorm:"foreignKey:UserID"`
	Dishes           []Dish    `gorm:"foreignKey:UserID"`
	Images           []Image   `gorm:"foreignKey:UserID"`
	SentMessages     []Message `gorm:"foreignKey:SenderID"`
	ReceivedMessages []Message `gorm:"foreignKey:RecipientID"`
	Todo             []Todo    `gorm:"foreignKey:UserID"`
	PrivilegeID      uint      `gorm:"not null"`
	Privilege        Privilege `gorm:"foreignKey:PrivilegeID"`
}

type Todo struct {
	gorm.Model
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	UserID      uint   `json:"userId"`
	User        User   `gorm:"foreignKey:UserID"`
}

// Privilege represents the user privilege levels
// e.g., Super Administrator, Administrator
// Slug: su, admin
// Users: relation to User

type Privilege struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Slug  string `gorm:"unique;not null"`
	Users []User `gorm:"foreignKey:PrivilegeID"`
}
