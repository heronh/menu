package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Dish represents a dish
type Dish struct {
	gorm.Model
	Name               string         `gorm:"not null"`
	Active             bool           `gorm:"default:true"`
	Description        string         `gorm:"not null"`
	Price              float64        `gorm:"not null"`
	ShowPrice          bool           `gorm:"default:true"`
	Time               pq.StringArray `gorm:"type:text[]"` // List with 3 times
	DaysOfWeekServed   pq.StringArray `gorm:"type:text[]"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	AuthorLastChangeID uint    `gorm:"not null"`
	AuthorLastChange   User    `gorm:"foreignKey:AuthorLastChangeID"`
	SectionID          uint    `gorm:"not null"`
	Section            Section `gorm:"foreignKey:SectionID"`
	CompanyID          uint    `gorm:"not null"`
	Company            Company `gorm:"foreignKey:CompanyID"`
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

// Section represents a dish section
type Section struct {
	gorm.Model
	Description string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorID    uint   `gorm:"not null"`
	Author      User   `gorm:"foreignKey:AuthorID"`
	Dishes      []Dish `gorm:"foreignKey:CategoryID"`
}
