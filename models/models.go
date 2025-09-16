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
