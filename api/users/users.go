package users

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	// Validate auth
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Initialize DB
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Get users
	var users []shared.User
	if err := shared.GetDB().
		Select("id, name, email, role, company_id, active, created_at").
		Where("active = ?", true). // Only active users
		Order("created_at DESC").  // Most recent first
		Limit(100).                // Prevent huge responses
		Find(&users).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=300")
	shared.SuccessResponse(w, users)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Validate auth (Check if admin?)
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var user shared.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// Hash password
	// Assuming user struct has Password field (temporary) or we use a separate struct.
	// Since shared.User likely has PasswordHash, we need to handle the plain password.
	// For now, let's assume the frontend sends 'password' which we need to hash.
	// But shared.User maps 'password_hash'. We might need a DTO.
	// Let's decode into a struct that has Password string.
	type CreateUserRequest struct {
		shared.User
		Password string `json:"password"`
	}
	var req CreateUserRequest
	// Re-decode to capture password (a bit inefficient but safe)
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		req.User.PasswordHash = string(hash)
	}

	if err := shared.GetDB().Create(&req.User).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed")
		return
	}

	shared.SuccessResponse(w, req.User)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	var user shared.User
	if err := shared.GetDB().Select("id, name, email, role, company_id, active, created_at").First(&user, id).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}
	shared.SuccessResponse(w, user)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	var user shared.User
	if err := shared.GetDB().First(&user, id).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	type UpdateUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Password string `json:"password,omitempty"`
	}
	var updates UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	user.Name = updates.Name
	user.Email = updates.Email
	user.Role = updates.Role

	if updates.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(updates.Password), bcrypt.DefaultCost)
		user.PasswordHash = string(hash)
	}

	if err := shared.GetDB().Save(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, user)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	if err := shared.GetDB().Model(&shared.User{}).Where("id = ?", id).Update("active", false).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "User deleted"})
}
func BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}
	if err := shared.GetDB().Model(&shared.User{}).Where("id = ?", id).Update("active", false).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}
	shared.SuccessResponse(w, map[string]interface{}{"success": true})
}

func AdminResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}
	// Reset to a default password (e.g., '123456') and require change
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err := shared.GetDB().Model(&shared.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password_hash":        string(hash),
		"must_change_password": true,
	}).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Reset failed")
		return
	}
	shared.SuccessResponse(w, map[string]interface{}{"success": true, "message": "Password reset to 123456"})
}
