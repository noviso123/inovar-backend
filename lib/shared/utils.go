package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT token for a user
func GenerateToken(user *User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not configured")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.ID, // Now a string (UUID)
		"email":     user.Email,
		"role":      user.Role,
		"companyId": user.CompanyID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token and returns user ID
func ValidateToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("invalid userID in claims")
	}
	return userID, nil
}

// GetAuthToken extracts token from Authorization header
func GetAuthToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// JSONResponse sends a JSON response
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse sends an error JSON response
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]interface{}{
		"success": false,
		"message": message,
	})
}

// SuccessResponse sends a success JSON response
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// ParseDateTime parses common ISO date strings
func ParseDateTime(s string) (time.Time, error) {
	formats := []string{time.RFC3339, "2006-01-02T15:04:05.000Z", "2006-01-02"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}
