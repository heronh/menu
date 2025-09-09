package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Privilege represents the user privilege levels
type Privilege struct {
	gorm.Model
	Name  string `gorm:"unique;not null"` // e.g., Super Administrator, Administrator
	Slug  string `gorm:"unique;not null"` // e.g., su, admin
	Users []User `gorm:"foreignKey:PrivilegeID"`
}

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
	PrivilegeID        uint       `gorm:"not null"`
	Privilege          Privilege  `gorm:"foreignKey:PrivilegeID"`
	AuthoredCategories []Category `gorm:"foreignKey:AuthorID"`
	AuthoredDishes     []Dish     `gorm:"foreignKey:AuthorLastChangeID"`
	InsertedImages     []Image    `gorm:"foreignKey:InsertedByID"`
	SentMessages       []Message  `gorm:"foreignKey:SenderID"`
	ReceivedMessages   []Message  `gorm:"foreignKey:RecipientID"`
	Tasks              []Task     `gorm:"foreignKey:UserID"`
}

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
	Users        []User    `gorm:"foreignKey:CompanyID"`
	Dishes       []Dish    `gorm:"foreignKey:CompanyID"`
	Images       []Image   `gorm:"foreignKey:CompanyID"`
	Messages     []Message `gorm:"foreignKey:SenderCompanyID"`
}

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

// Log represents a log entry for model changes
type Log struct {
	gorm.Model
	Text           string `gorm:"not null"`
	CreatedAt      time.Time
	ModelChangedID *uint  // Pointer to allow null, representing the ID of the changed record
	ModelType      string // e.g., "User", "Company", "Dish" - to know which table ModelChangedID refers to
}

// Task represents a task for a user
type Task struct {
	gorm.Model
	Text      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint `gorm:"not null"`
	User      User `gorm:"foreignKey:UserID"`
	Finished  bool `gorm:"default:false"`
}

type Todo struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"primary_key"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Description string    `json:"description"`
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
