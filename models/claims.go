package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"username" db:"username"`
	UserRole string `json:"userRole" db:"role"`
	UserID   string `json:"userID" db:"user_id"`
	jwt.StandardClaims
}
