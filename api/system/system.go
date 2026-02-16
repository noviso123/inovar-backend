package system

import (
	"encoding/json"
	"fmt"
	"inovar/lib/shared"
	"net/http"
	"os"
)

func WhatsAppStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Mock status
	shared.SuccessResponse(w, map[string]interface{}{
		"enabled":   true,
		"connected": false,
		"qrCode":    "mock_qr_code_base64",
	})
}

func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Mock routes list or dynamic extraction
	shared.SuccessResponse(w, []string{
		"/api/users", "/api/clients", "/api/requests", "/api/equipments",
		"/api/agenda", "/api/company", "/api/finance", "/api/audit",
		"/api/settings", "/api/fiscal", "/api/system",
	})
}

func TablesHandler(w http.ResponseWriter, r *http.Request) {
	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var tables []string
	// Query real tables from postgres
	err := shared.GetDB().Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables).Error
	if err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed: "+err.Error())
		return
	}
	shared.SuccessResponse(w, tables)
}

func BucketsHandler(w http.ResponseWriter, r *http.Request) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Supabase credentials not configured")
		return
	}

	url := fmt.Sprintf("%s/storage/v1/bucket", supabaseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Request failed: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	shared.SuccessResponse(w, data)
}

func TableDataHandler(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	if tableName == "" {
		shared.ErrorResponse(w, http.StatusBadRequest, "Table name required")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var data []map[string]interface{}
	if err := shared.GetDB().Table(tableName).Limit(100).Find(&data).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed: "+err.Error())
		return
	}

	shared.SuccessResponse(w, data)
}
