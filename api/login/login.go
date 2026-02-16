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
		Where("email = ?", req.Email). // Case sensitivity depends on DB, but email should be unique
		First(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	// Check if active
	if !user.Active {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Sua conta está inativa. Entre em contato com o suporte.")
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

// ChangePasswordRequest struct
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := shared.GetAuthToken(r)
	userID, err := shared.ValidateToken(token)
	if err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var user shared.User
	if err := shared.GetDB().First(&user, userID).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Senha atual incorreta")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hashed)
	user.MustChangePassword = false // Reset flag

	if err := shared.GetDB().Save(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Senha alterada com sucesso"})
}

// ForgotPasswordRequest struct
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Just a stub for now - requires email service
	// We verify if user exists but don't reveal it (security best practice) causes ambiguity,
	// but here we can return success always.
	// For MVP, we can just say "If email exists, we sent instructions"

	shared.SuccessResponse(w, map[string]string{"message": "Se o e-mail existir, enviaremos instruções."})
}

// ResetPasswordRequest struct
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Stub - requires token verification logic
	shared.SuccessResponse(w, map[string]string{"message": "Senha redefinida com sucesso (Stub)"})
}

// UpdateProfileRequest - simplified version of user update
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := shared.GetAuthToken(r)
	userID, err := shared.ValidateToken(token)
	if err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var user shared.User
	if err := shared.GetDB().First(&user, userID).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	// Update allowed fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := shared.GetDB().Save(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	shared.SuccessResponse(w, user)
}
