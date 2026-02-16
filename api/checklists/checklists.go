package checklists

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
	"strconv"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	token := shared.GetAuthToken(r)
	if token == "" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if _, err := shared.ValidateToken(token); err != nil {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestId := r.PathValue("requestId")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var items []shared.ChecklistItem
	if err := shared.GetDB().Where("request_id = ?", requestId).Find(&items).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	shared.SuccessResponse(w, items)
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

	requestId := r.PathValue("requestId")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var item shared.ChecklistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// Force RequestID from path
	// Convert string path value to uint
	if id, err := strconv.Atoi(requestId); err == nil {
		item.RequestID = uint(id)
	}

	if err := shared.GetDB().Create(&item).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed")
		return
	}

	shared.SuccessResponse(w, item)
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

	itemId := r.PathValue("itemId")

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var item shared.ChecklistItem
	if err := shared.GetDB().First(&item, itemId).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "Item not found")
		return
	}

	var updates shared.ChecklistItem
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	if err := shared.GetDB().Model(&item).Updates(updates).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, item)
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

	itemId := r.PathValue("itemId")
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Delete(&shared.ChecklistItem{}, itemId).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Item deleted"})
}
