package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	defaultSecret = []byte("your-secret-key-change-in-production")
)

// Claims represents JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GetJWTSecret returns the JWT secret key
func GetJWTSecret() []byte {
	if jwtSecret != nil {
		return jwtSecret
	}

	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		jwtSecret = []byte(secret)
		return jwtSecret
	}

	// Use default secret for development (should be changed in production)
	jwtSecret = defaultSecret
	return jwtSecret
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user *User) (string, error) {
	expirationHours := 24 // Default 24 hours
	if hours := os.Getenv("JWT_EXPIRATION_HOURS"); hours != "" {
		// Parse hours from env (simplified, could add proper parsing)
		// For now, keep default
	}

	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-dashboard",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
