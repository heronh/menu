package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name      string
	Email     string `gorm:"uniqueIndex"`
	Password  string // Hashed password
	RoleID    uint
	Role      Role
	CompanyID uint
	Company   Company
}

// Company represents a company
type Company struct {
	gorm.Model
	Name  string
	Users []User
}

// Role represents a user role
type Role struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex"` // e.g., admin, editor, viewer
	Users []User
}

// Dish represents a menu item
type Dish struct {
	gorm.Model
	Name         string
	Description  string
	Price        float64
	CategoryID   uint
	Category     DishCategory
	ImageID      uint
	Image        Image
	RestaurantID uint // Assuming dishes belong to a company (restaurant)
}

// DishCategory represents a category for dishes
type DishCategory struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex"`
	Dishes []Dish
}

// Image represents an image, could be for a dish or other entities
type Image struct {
	gorm.Model
	URL        string `gorm:"uniqueIndex"` // URL of the image
	AltText    string
	UploadedAt time.Time
}
