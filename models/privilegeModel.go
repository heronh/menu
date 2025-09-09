package models

import "gorm.io/gorm"

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
