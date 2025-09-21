package models

import (
	"time"

	"gorm.io/gorm"
)

// Log represents a log entry for model changes
type Log struct {
	gorm.Model
	Text           string `gorm:"not null"`
	CreatedAt      time.Time
	ModelChangedID *uint  // Pointer to allow null, representing the ID of the changed record
	ModelType      string // e.g., "User", "Company", "Dish" - to know which table ModelChangedID refers to
}
