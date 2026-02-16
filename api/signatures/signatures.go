package signatures

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
	"strconv"
)

// SaveHandler - POST /api/requests/{id}/assinatura
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	requestIdStr := r.PathValue("id")
	requestID, err := strconv.ParseUint(requestIdStr, 10, 32)
	if err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid Request ID")
		return
	}

	var req shared.Assinatura
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.RequestID = uint(requestID)

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Create(&req).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to save signature")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
