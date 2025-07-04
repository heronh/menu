package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecretKey is used to sign and verify JWT tokens.
// It's critical to keep this secret and preferably load it from environment variables.
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

const GinContextKey = "userClaims"

// Initialize loads the JWT_SECRET_KEY from environment variables.
// Call this function at application startup.
func Initialize() {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		// For development, a default key can be used, but this is insecure for production.
		// In a real application, you might want to panic or log a fatal error if the key isn't set.
		fmt.Println("Warning: JWT_SECRET_KEY environment variable not set. Using default insecure key.")
		jwtSecretKey = []byte("a_very_insecure_default_secret_key_replace_me")
	} else {
		jwtSecretKey = []byte(secret)
	}
}

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plain text password with a bcrypt hashed password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT creates a new JWT for a given user.
func GenerateJWT(user *models.User, privilegeName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Privilege: privilegeName,
		CompanyID: user.CompanyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// JWTMiddleware creates a Gin middleware for JWT authentication.
// It verifies the token from the "Authorization" header.
// If the token is valid, it stores the user claims in the Gin context.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Store claims in context for use by handlers
		c.Set(GinContextKey, claims)
		c.Next()
	}
}

// GetClaimsFromContext retrieves the user claims from the Gin context.
// Returns nil if claims are not found (e.g., middleware not used or token invalid).
func GetClaimsFromContext(c *gin.Context) *Claims {
	claims, exists := c.Get(GinContextKey)
	if !exists {
		return nil
	}
	if castedClaims, ok := claims.(*Claims); ok {
		return castedClaims
	}
	return nil
}

// Authorize middleware checks if the user has one of the required privileges.
func Authorize(requiredPrivileges ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaimsFromContext(c)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied. User not authenticated properly."})
			return
		}

		userPrivilege := claims.Privilege
		allowed := false
		for _, requiredPrivilege := range requiredPrivileges {
			if userPrivilege == requiredPrivilege {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Access denied. Required privilege: %s", strings.Join(requiredPrivileges, " or "))})
			return
		}

		c.Next()
	}
}
