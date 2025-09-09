package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Category represents a dish category
type Category struct {
	gorm.Model
	Description string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorID    uint   `gorm:"not null"`
	Author      User   `gorm:"foreignKey:AuthorID"`
	Dishes      []Dish `gorm:"foreignKey:CategoryID"`
}

// Dish represents a dish
type Dish struct {
	gorm.Model
	Description        string         `gorm:"not null"`
	Value              float64        `gorm:"not null"`
	Time               pq.StringArray `gorm:"type:text[]"` // List with 3 times
	DaysOfWeekServed   pq.StringArray `gorm:"type:text[]"`
	Active             bool           `gorm:"default:true"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	AuthorLastChangeID uint     `gorm:"not null"`
	AuthorLastChange   User     `gorm:"foreignKey:AuthorLastChangeID"`
	CategoryID         uint     `gorm:"not null"`
	Category           Category `gorm:"foreignKey:CategoryID"`
	CompanyID          uint     `gorm:"not null"`
	Company            Company  `gorm:"foreignKey:CompanyID"`
}

// Image represents an image file
type Image struct {
	gorm.Model
	OriginalFileName string  `gorm:"not null"`
	UniqueName       string  `gorm:"unique;not null"`
	Storage          string  // e.g., local, s3
	InsertedByID     uint    `gorm:"not null"`
	InsertedBy       User    `gorm:"foreignKey:InsertedByID"`
	CompanyID        uint    `gorm:"not null"`
	Company          Company `gorm:"foreignKey:CompanyID"`
	CreatedAt        time.Time
}

// Log represents a log entry for model changes
type Log struct {
	gorm.Model
	Text           string `gorm:"not null"`
	CreatedAt      time.Time
	ModelChangedID *uint  // Pointer to allow null, representing the ID of the changed record
	ModelType      string // e.g., "User", "Company", "Dish" - to know which table ModelChangedID refers to
}

// Ensure all models have CreatedAt and UpdatedAt managed by GORM by default
// For models that have explicit CreatedAt and UpdatedAt, GORM handles them.
// For gorm.Model, it includes ID, CreatedAt, UpdatedAt, DeletedAt.
// I've added CreatedAt and UpdatedAt to models that didn't explicitly list them
// but are typical to have, assuming GORM's default behavior is desired unless specified.
// For User, Company, Category, Dish, Task, explicit CreatedAt and UpdatedAt are defined.
// For Image, Message, Log, CreatedAt is defined or implied by gorm.Model.
// Added Slug to Privilege for easier referencing (e.g. "su" for Super Administrator).
// Used pq.StringArray for Time and DaysOfWeekServed in Dish model. This requires `github.com/lib/pq`.
// For foreign keys that can be null (User.CompanyID), used pointer type (*uint).
// Clarified Log.ModelChangedID and added ModelType to specify which table the ID refers to.
// Added relationships (e.g. Privilege.Users, User.AuthoredCategories) for easier GORM operations.
