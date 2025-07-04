package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct for JWT, including user details
type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Privilege string `json:"privilege"`
	CompanyID uint   `json:"company_id"`
	jwt.RegisteredClaims
}
