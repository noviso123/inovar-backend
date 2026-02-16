package company

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
)

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

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Assuming single company for now, or get by ID if SaaS.
	// Model has CompanyID. Let's get the first one or the one linked to user.
	var company shared.Company
	if err := shared.GetDB().First(&company).Error; err != nil {
		// Return empty or error?
		shared.ErrorResponse(w, http.StatusNotFound, "Company not found")
		return
	}

	shared.SuccessResponse(w, company)
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

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var company shared.Company
	if err := shared.GetDB().First(&company).Error; err != nil {
		// Create if not exists?
		company = shared.Company{}
	}

	var updates shared.Company
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// If no ID, create
	if company.ID == 0 {
		if err := shared.GetDB().Create(&updates).Error; err != nil {
			shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed")
			return
		}
		shared.SuccessResponse(w, updates)
		return
	}

	if err := shared.GetDB().Model(&company).Updates(updates).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Update failed")
		return
	}

	shared.SuccessResponse(w, company)
}
