package auth

import (
	"testing"
	"time"

	"example.com/m/v2/models"
	"github.com/golang-jwt/jwt/v5"
)

// TestHashAndCheckPassword verifies that password hashing and checking work correctly.
func TestHashAndCheckPassword(t *testing.T) {
	password := "plainTextPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hashedPassword == password {
		t.Errorf("Hashed password should not be the same as plain text password.")
	}

	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("CheckPasswordHash() failed to verify correct password.")
	}

	if CheckPasswordHash("wrongPassword", hashedPassword) {
		t.Errorf("CheckPasswordHash() incorrectly verified a wrong password.")
	}
}

// TestGenerateAndParseJWT is a basic test for JWT generation and parsing.
// Note: For more robust JWT testing, especially involving time-based claims,
// you might need to mock time or use libraries that allow time manipulation.
func TestGenerateAndParseJWT(t *testing.T) {
	// Initialize must be called to set jwtSecretKey, especially if it relies on env vars.
	// For testing, we can set a fixed secret or ensure Initialize uses a predictable test secret.
	originalKey := jwtSecretKey
	jwtSecretKey = []byte("test_secret_key_for_auth_test") // Override for test
	defer func() { jwtSecretKey = originalKey }()          // Restore original key

	user := &models.User{
		Model:     gorm.Model{ID: 1},
		Email:     "test@example.com",
		CompanyID: 10,
	}
	privilegeName := "Manager"

	tokenString, err := GenerateJWT(user, privilegeName)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	if tokenString == "" {
		t.Fatal("Generated JWT string is empty.")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		t.Fatalf("ParseWithClaims() error = %v", err)
	}

	if !token.Valid {
		t.Errorf("Parsed token is not valid.")
	}

	if claims.UserID != user.ID {
		t.Errorf("Parsed claims UserID got = %v, want = %v", claims.UserID, user.ID)
	}
	if claims.Email != user.Email {
		t.Errorf("Parsed claims Email got = %v, want = %v", claims.Email, user.Email)
	}
	if claims.Privilege != privilegeName {
		t.Errorf("Parsed claims Privilege got = %v, want = %v", claims.Privilege, privilegeName)
	}
	if claims.CompanyID != user.CompanyID {
		t.Errorf("Parsed claims CompanyID got = %v, want = %v", claims.CompanyID, user.CompanyID)
	}

	// Check expiry (approximate, as it's set to Now + duration)
	// This is a very basic check. More precise checks might involve mocking time.Now().
	expectedExpiry := time.Now().Add(24 * time.Hour)
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) || claims.ExpiresAt.Time.After(expectedExpiry.Add(1*time.Minute)) {
		// Allow a small window for execution time
		t.Errorf("Token expiry time seems incorrect. Got: %v, Expected around: %v", claims.ExpiresAt.Time, expectedExpiry)
	}
}

// Placeholder for GORM model to satisfy User model's gorm.Model field in tests
// In a real GORM test setup, you'd typically use a test database.
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Ensure models.User's gorm.Model is compatible, or use a mock if not interacting with DB.
// For these specific tests, User struct is simple enough.
// We need to mock gorm.Model for models.User
type gorm struct {
	Model models.Model // This is a placeholder type
}
