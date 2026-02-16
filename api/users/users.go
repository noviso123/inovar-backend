package users

import (
	"crypto/rand"
	"encoding/hex"
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
		Select("id, name, email, role, company_id, active, phone, avatar_url, created_at").
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

// CreateUserRequest combines User fields with a plain password
type CreateUserRequest struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Phone    *string `json:"phone"`
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Decode body ONCE into a combined struct
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// Create the user model
	user := shared.User{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
		Phone: req.Phone,
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = generateUUID()
	}

	// Hash password
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.PasswordHash = string(hash)
	}

	// Set defaults
	active := true
	user.Active = &active
	mustChange := true
	user.MustChangePassword = &mustChange

	if err := shared.GetDB().Create(&user).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed: "+err.Error())
		return
	}

	shared.SuccessResponse(w, user)
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
	if err := shared.GetDB().
		Select("id, name, email, role, company_id, active, phone, avatar_url, created_at").
		Where("id = ?", id).
		First(&user).Error; err != nil {
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
	if err := shared.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
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

	if updates.Name != "" {
		user.Name = updates.Name
	}
	if updates.Email != "" {
		user.Email = updates.Email
	}
	if updates.Role != "" {
		user.Role = updates.Role
	}

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
	// Auth check
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

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
	// Auth check
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

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
