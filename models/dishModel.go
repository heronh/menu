package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Dish represents a dish
type Dish struct {
	gorm.Model
	Name            string         `gorm:"not null"`
	Active          bool           `gorm:"default:true"`
	Description     string         `gorm:"type:text"`
	ShowDescription bool           `gorm:"default:true"`
	Price           float64        `gorm:"not null"`
	ShowPrice       bool           `gorm:"default:true"`
	Availability    pq.StringArray `gorm:"type:text[]"`
	WeekDays        pq.StringArray `gorm:"type:text[]"`
	UserID          uint           `gorm:"not null"`
	User            User           `gorm:"foreignKey:UserID"`
	SectionID       uint           `gorm:"not null"`
	Section         Section        `gorm:"foreignKey:SectionID"`
	CompanyID       uint           `gorm:"not null"`
	Company         Company        `gorm:"foreignKey:CompanyID"`
	Images          []Image        `gorm:"many2many:dish_images;"`
	AvailableNow    bool           `gorm:"-:all"` // Transient field, not stored in DB
}

// Image represents an image file
type Image struct {
	gorm.Model
	OriginalFileName string  `gorm:"not null"`
	UniqueName       string  `gorm:"unique;not null"`
	Storage          string  `gorm:"default:'local'"`
	UserID           uint    `gorm:"not null"`
	User             User    `gorm:"foreignKey:UserID"`
	CompanyID        uint    `gorm:"not null"`
	Company          Company `gorm:"foreignKey:CompanyID"`
	CreatedAt        time.Time
}

// Section represents a dish section
type Section struct {
	gorm.Model
	Description string `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	Dishes      []Dish
	CompanyID   uint    `gorm:"not null"`
	Company     Company `gorm:"foreignKey:CompanyID"`
}
