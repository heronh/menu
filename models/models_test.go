package models

import (
	"testing"
	"time"

	// "gorm.io/driver/sqlite" // Example for a test DB
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing GORM models.
// Returns a pointer to gorm.DB and a cleanup function.
// This is a common pattern for model testing.
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// For actual tests, use an in-memory SQLite or a dedicated test PostgreSQL instance.
	// db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	// if err != nil {
	// 	t.Fatalf("Failed to connect to test database: %v", err)
	// }

	// For this placeholder, we'll return nil.
	// Replace with actual DB setup when implementing tests.
	db := (*gorm.DB)(nil)
	if db == nil {
		t.Skip("Skipping model tests: Test DB not configured.")
		return nil, func() {}
	}

	// Run migrations for all models
	err := db.AutoMigrate(
		&Privilege{}, &User{}, &Company{}, &Category{},
		&Dish{}, &Image{}, &Message{}, &Log{}, &Todo{},
	)
	if err != nil {
		t.Fatalf("Failed to auto-migrate test database: %v", err)
	}

	cleanup := func() {
		// Clean up: drop tables or close connection if necessary
		// For in-memory SQLite, closing might not be needed or tables are dropped automatically.
		// If using a persistent test DB, you'd drop tables:
		// sqlDB, _ := db.DB()
		// sqlDB.Close()
		// Or db.Migrator().DropTable(...) for all tables
	}

	return db, cleanup
}

// TestUserHooks checks BeforeCreate and BeforeUpdate hooks for User model.
func TestUserHooks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	} // Skip if DB setup was skipped
	defer cleanup()

	now := time.Now()
	user := User{Name: "Test User Hook", Email: "hook@example.com"}

	// Test BeforeCreate
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if user.CreationDate.IsZero() || user.CreationDate.Before(now.Add(-1*time.Second)) {
		t.Errorf("User CreationDate not set correctly by BeforeCreate hook. Got: %v", user.CreationDate)
	}
	if user.LastModified.IsZero() || !user.LastModified.Equal(user.CreationDate) {
		t.Errorf("User LastModified not set correctly by BeforeCreate hook. Got: %v, Expected: %v", user.LastModified, user.CreationDate)
	}

	// Test BeforeUpdate
	oldLastModified := user.LastModified
	time.Sleep(10 * time.Millisecond) // Ensure time changes
	user.Name = "Updated Test User Hook"
	if err := db.Save(&user).Error; err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	if user.LastModified.Equal(oldLastModified) || user.LastModified.Before(oldLastModified) {
		t.Errorf("User LastModified not updated by BeforeUpdate hook. Got: %v, Old: %v", user.LastModified, oldLastModified)
	}
}

// TestCompanyHooks (Similar structure for Company)
func TestCompanyHooks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()
	// ... implementation similar to TestUserHooks ...
	t.Log("Company hooks test placeholder - implement similarly to User hooks.")
}

// TestDishTimeAndDaysOfWeekServed checks GORM's handling of []string for text[] in PostgreSQL.
// This test would ideally run against a PostgreSQL test instance.
// With SQLite, []string might be stored as JSON or a delimited string, depending on GORM dialect.
func TestDishStringArrayFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	// This test is more meaningful with PostgreSQL.
	// If using SQLite for tests, GORM might serialize []string differently.
	// For now, just a basic create and read.

	// Need a Category and Company for the Dish
	authorUser := User{Name: "Dish Author", Email: "dishauthor@example.com"}
	db.Create(&authorUser)

	company := Company{Name: "Dish Test Co", CNPJ: "DISHCO123"}
	db.Create(&company)

	category := Category{Description: "Test Cat for Dish", AuthorID: authorUser.ID}
	db.Create(&category)

	dishTimes := []string{"10:00", "14:00", "18:00"}
	dishDays := []string{"Mon", "Wed", "Fri"}

	dish := Dish{
		Description:                "Array Test Dish",
		Value:                      19.99,
		Time:                       dishTimes,
		DaysOfWeekServed:           dishDays,
		CategoryID:                 category.ID,
		CompanyID:                  company.ID,
		AuthorOfLastModificationID: authorUser.ID,
	}

	if err := db.Create(&dish).Error; err != nil {
		t.Fatalf("Failed to create dish: %v", err)
	}

	var fetchedDish Dish
	if err := db.First(&fetchedDish, dish.ID).Error; err != nil {
		t.Fatalf("Failed to fetch dish: %v", err)
	}

	if len(fetchedDish.Time) != len(dishTimes) {
		t.Errorf("Fetched dish times length mismatch. Got %d, want %d. (Times: %v)", len(fetchedDish.Time), len(dishTimes), fetchedDish.Time)
	} else {
		for i, tm := range dishTimes {
			if fetchedDish.Time[i] != tm {
				t.Errorf("Fetched dish time at index %d mismatch. Got %s, want %s", i, fetchedDish.Time[i], tm)
			}
		}
	}

	if len(fetchedDish.DaysOfWeekServed) != len(dishDays) {
		t.Errorf("Fetched dish days length mismatch. Got %d, want %d. (Days: %v)", len(fetchedDish.DaysOfWeekServed), len(dishDays), fetchedDish.DaysOfWeekServed)
	} else {
		for i, day := range dishDays {
			if fetchedDish.DaysOfWeekServed[i] != day {
				t.Errorf("Fetched dish day at index %d mismatch. Got %s, want %s", i, fetchedDish.DaysOfWeekServed[i], day)
			}
		}
	}
	t.Logf("Note: []string field test is more effective with a PostgreSQL test database. SQLite behavior might differ.")
}

// Add more model-specific tests here, e.g.:
// - Testing unique constraints (requires attempting to create duplicates).
// - Testing foreign key relationships (creating related records, querying them).
// - Testing soft delete (gorm.DeletedAt).
// - Testing custom validation logic if added to models.
// - Testing specific data type handling if complex types are used.
