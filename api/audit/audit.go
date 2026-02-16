package audit

import (
	"inovar/lib/shared"
	"net/http"
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

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	query := shared.GetDB().Model(&shared.AuditLog{})

	entity := r.URL.Query().Get("entity")
	if entity != "" {
		query = query.Where("entity = ?", entity)
	}

	userId := r.URL.Query().Get("userId")
	if userId != "" {
		query = query.Where("user_id = ?", userId)
	}

	var logs []shared.AuditLog
	if err := query.Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	shared.SuccessResponse(w, logs)
}
