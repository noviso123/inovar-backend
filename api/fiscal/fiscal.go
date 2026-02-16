package fiscal

import (
	"inovar/lib/shared"
	"net/http"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Return mock config or from DB
	shared.SuccessResponse(w, map[string]interface{}{
		"environment":      "homologacao",
		"regimeTributario": "simples_nacional",
	})
}

func CertificateHandler(w http.ResponseWriter, r *http.Request) {
	// Mock upload
	shared.SuccessResponse(w, map[string]string{"message": "Certificate uploaded"})
}
