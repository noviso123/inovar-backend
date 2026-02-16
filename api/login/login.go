package login

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// CORS - allow frontend domain
	origin := r.Header.Get("Origin")
	if origin == "https://inovar-gestao.vercel.app" || origin == "http://localhost:5173" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Initialize DB
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Parse request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Find user
	var user shared.User
	if err := shared.GetDB().
		Select("id, name, email, password_hash, role, company_id, must_change_password, active").
		Where("email = ? AND active = ?", req.Email, true).
		First(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	// Generate token
	token, err := shared.GenerateToken(&user)
	if err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Token generation failed")
		return
	}

	// Send response
	shared.SuccessResponse(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id":                 user.ID,
			"name":               user.Name,
			"email":              user.Email,
			"role":               user.Role,
			"companyId":          user.CompanyID,
			"mustChangePassword": user.MustChangePassword,
		},
		"accessToken":  token,
		"refreshToken": token, // Same for now
		"expiresIn":    86400, // 24h
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	// CORS - allow frontend domain
	origin := r.Header.Get("Origin")
	if origin == "https://inovar-gestao.vercel.app" || origin == "http://localhost:5173" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Auth required
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := shared.ValidateToken(token)
	if err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var user shared.User
	if err := shared.GetDB().
		Select("id, name, email, role, company_id, must_change_password").
		First(&user, userID).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	shared.SuccessResponse(w, user)
}
