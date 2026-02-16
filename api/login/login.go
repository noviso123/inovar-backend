package login

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"inovar/lib/shared"
	"log"
	"net/http"
	"time"

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

	// Find user - fetch ALL columns needed for auth (no hardcoded Select)
	var user shared.User
	if err := shared.GetDB().
		Where("email = ?", req.Email).
		First(&user).Error; err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err)
		shared.ErrorResponse(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	// Check if active
	if user.Active != nil && !*user.Active {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Sua conta está inativa.")
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
			"active":             user.Active,
			"mustChangePassword": user.MustChangePassword,
			"avatarUrl":          user.AvatarURL,
			"phone":              user.Phone,
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
		Where("id = ?", userID).
		First(&user).Error; err != nil {
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
	if err := shared.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Senha atual incorreta")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hashed)
	f := false
	user.MustChangePassword = &f // Reset flag

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

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Look up user by email
	var user shared.User
	if err := shared.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Security: Don't reveal if email exists or not
		log.Printf("Forgot password: email %s not found", req.Email)
		shared.SuccessResponse(w, map[string]string{"message": "Se o e-mail existir, enviaremos instruções."})
		return
	}

	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	resetToken := hex.EncodeToString(tokenBytes)

	// Delete any existing tokens for this user
	shared.GetDB().Where("user_id = ?", user.ID).Delete(&shared.PasswordResetToken{})

	// Store the token in the database with 1-hour expiration
	tokenRecord := shared.PasswordResetToken{
		ID:        generateUUID(),
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := shared.GetDB().Create(&tokenRecord).Error; err != nil {
		log.Printf("Failed to store reset token: %v", err)
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// Log the reset link (MVP - no email service yet)
	log.Printf("🔑 PASSWORD RESET TOKEN for %s: %s", req.Email, resetToken)
	log.Printf("🔗 Reset link: /reset-password?token=%s", resetToken)

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

	if req.Token == "" {
		shared.ErrorResponse(w, http.StatusBadRequest, "Token é obrigatório")
		return
	}

	if len(req.NewPassword) < 6 {
		shared.ErrorResponse(w, http.StatusBadRequest, "A senha deve ter no mínimo 6 caracteres")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Find the token record
	var tokenRecord shared.PasswordResetToken
	if err := shared.GetDB().
		Where("token = ? AND used = ?", req.Token, false).
		First(&tokenRecord).Error; err != nil {
		log.Printf("Reset password: invalid token")
		shared.ErrorResponse(w, http.StatusBadRequest, "Token inválido ou já utilizado")
		return
	}

	// Check expiration
	if time.Now().After(tokenRecord.ExpiresAt) {
		log.Printf("Reset password: token expired for user %s", tokenRecord.UserID)
		shared.ErrorResponse(w, http.StatusBadRequest, "Token expirado. Solicite um novo link de recuperação.")
		return
	}

	// Hash the new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update the user's password
	if err := shared.GetDB().Model(&shared.User{}).
		Where("id = ?", tokenRecord.UserID).
		Updates(map[string]interface{}{
			"password_hash":        string(hashed),
			"must_change_password": false,
		}).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Mark token as used
	shared.GetDB().Model(&tokenRecord).Update("used", true)

	log.Printf("✅ Password reset successful for user %s", tokenRecord.UserID)
	shared.SuccessResponse(w, map[string]string{"message": "Senha redefinida com sucesso"})
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
	if err := shared.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
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

// LogoutHandler - POST /api/auth/logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// For JWT, logout is primarily handled by the client clearing the token.
	// We return success to satisfy the frontend.
	shared.SuccessResponse(w, map[string]string{"message": "Logged out successfully"})
}

// generateUUID creates a simple UUID v4
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	return hex.EncodeToString(b[:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}
