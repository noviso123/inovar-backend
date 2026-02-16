package system

import (
	"inovar/lib/shared"
	"net/http"
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
	// Mock routes list
	shared.SuccessResponse(w, []string{"/api/users", "/api/clients", "/api/requests"})
}
