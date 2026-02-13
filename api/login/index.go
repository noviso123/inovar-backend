package handler

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

func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS - allow frontend domain and local development
	origin := r.Header.Get("Origin")
	allowedOrigins := map[string]bool{
		"https://inovar-gestao.vercel.app": true,
		"http://localhost:5173":            true,
		"http://localhost:3000":            true,
		"http://localhost:3001":            true,
	}

	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		// Fallback for other environments if needed, or strict deny
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 初始化 DB antes de processar
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Find user - optimized query with index (email) and minimal fields
	var user shared.User
	if err := shared.GetDB().
		Select("id, name, email, password_hash, role, company_id, must_change_password, active").
		Where("email = ? AND active = ?", req.Email, true). // Use index + filter active
		First(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Credenciais inv├ílidas")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Credenciais inv├ílidas")
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
