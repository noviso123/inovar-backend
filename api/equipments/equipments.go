package equipments

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

	// Filters
	query := shared.GetDB().Model(&shared.Equipment{})

	clientId := r.URL.Query().Get("clientId")
	if clientId != "" {
		query = query.Where("client_id = ?", clientId)
	}

	activeOnly := r.URL.Query().Get("activeOnly")
	if activeOnly == "true" {
		query = query.Where("active = ?", true)
	}

	var equipments []shared.Equipment
	if err := query.Find(&equipments).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	shared.SuccessResponse(w, equipments)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	var equip shared.Equipment
	if err := json.NewDecoder(r.Body).Decode(&equip); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	if err := shared.GetDB().Create(&equip).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed")
		return
	}

	shared.SuccessResponse(w, equip)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var equip shared.Equipment
	if err := shared.GetDB().First(&equip, id).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "Equipment not found")
		return
	}

	shared.SuccessResponse(w, equip)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var equip shared.Equipment
	if err := shared.GetDB().First(&equip, id).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "Equipment not found")
		return
	}

	var updates shared.Equipment
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	if err := shared.GetDB().Model(&equip).Updates(updates).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, equip)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Delete(&shared.Equipment{}, id).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Equipment deleted"})
}

func DeactivateHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Model(&shared.Equipment{}).Where("id = ?", id).Update("active", false).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, map[string]interface{}{"success": true})
}

func ReactivateHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Model(&shared.Equipment{}).Where("id = ?", id).Update("active", true).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, map[string]interface{}{"success": true})
}
