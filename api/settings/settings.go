package settings

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

	var settings []shared.Setting
	if err := shared.GetDB().Find(&settings).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	// Convert to map
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	shared.SuccessResponse(w, settingsMap)
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

	type SettingsRequest struct {
		Settings map[string]string `json:"settings"`
	}
	var req SettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// Update each setting
	for k, v := range req.Settings {
		// Upsert
		var setting shared.Setting
		if err := shared.GetDB().Where("key = ?", k).First(&setting).Error; err != nil {
			// Create
			shared.GetDB().Create(&shared.Setting{Key: k, Value: v})
		} else {
			// Update
			setting.Value = v
			shared.GetDB().Save(&setting)
		}
	}

	shared.SuccessResponse(w, map[string]string{"message": "Settings updated"})
}
