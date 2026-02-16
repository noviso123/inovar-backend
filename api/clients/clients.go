package clients

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	// Auth required
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// List clients - optimized query
	var clients []shared.Client
	if err := shared.GetDB().
		Where("active = ?", true).
		Order("name ASC"). // Alphabetical order
		Limit(200).        // Prevent huge responses
		Find(&clients).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	shared.SuccessResponse(w, clients)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Auth required
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Create client
	var client shared.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	if err := shared.GetDB().Create(&client).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed")
		return
	}

	shared.SuccessResponse(w, client)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// Auth required
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	var client shared.Client
	if err := shared.GetDB().Where("id = ?", id).First(&client).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "Client not found")
		return
	}

	shared.SuccessResponse(w, client)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Auth required
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	var client shared.Client
	if err := shared.GetDB().Where("id = ?", id).First(&client).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "Client not found")
		return
	}

	var updates shared.Client
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	if err := shared.GetDB().Model(&client).Updates(updates).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, client)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Auth required
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	// Soft delete (set active = false) or hard delete? 'active = false' is safer for relational integrity.
	// But let's check shared.Client model. Assuming it has DeletedAt (gorm.Model) or Active field.
	// The snippet showed "Where active = true", so we should set active = false.

	if err := shared.GetDB().Model(&shared.Client{}).Where("id = ?", id).Update("active", false).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Client deleted"})
}

// BlockHandler - PATCH /api/clients/{id}/block
func BlockHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	id := r.PathValue("id")
	if err := shared.GetDB().Model(&shared.Client{}).Where("id = ?", id).Update("active", false).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Block failed")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Client blocked"})
}
