package fiscal

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
	"strconv"
	"time"
)

// ConfigHandler (Existing)
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	// ... (Implementation already exists in original file? No, I need to preserve it or append if I'm rewriting)
	// Wait, I should not overwrite the whole file if it has content.
	// Use append/replace if possible?
	// The previous tool created `api/fiscal/fiscal.go` with Stubs.
	// I will rewrite the whole file to include both Config/Certificate AND NFSe handlers to be clean.

	if r.Method != "GET" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Stub
	shared.SuccessResponse(w, map[string]interface{}{
		"environment": "homologacao",
		"regime":      "simples_nacional",
	})
}

// CertificateHandler (Existing Stub)
func CertificateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	shared.SuccessResponse(w, map[string]string{"message": "Certificado enviado (Stub)"})
}

// NFSe Handlers

// IssueNFSeHandler - POST /api/requests/{id}/nfse
func IssueNFSeHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Mock Issue Logic
	nfse := shared.NFSe{
		RequestID: uint(requestID),
		Numero:    "2023000" + requestIdStr, // Mock number
		Status:    "emitido",
		PDFURL:    "https://example.com/nfse.pdf",
		XMLURL:    "https://example.com/nfse.xml",
		CreatedAt: time.Now(),
	}

	if err := shared.GetDB().Create(&nfse).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to issue NFSe")
		return
	}

	shared.SuccessResponse(w, nfse)
}

// GetNFSeHandler - GET /api/requests/{id}/nfse
func GetNFSeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	requestIdStr := r.PathValue("id")

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	var nfse shared.NFSe
	if err := shared.GetDB().Where("request_id = ?", requestIdStr).First(&nfse).Error; err != nil {
		shared.ErrorResponse(w, http.StatusNotFound, "NFSe not found")
		return
	}

	shared.SuccessResponse(w, nfse)
}

// CancelNFSeHandler - DELETE /api/requests/{id}/nfse
func CancelNFSeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Stub
	shared.SuccessResponse(w, map[string]string{"message": "NFSe cancelada (Stub)"})
}

// CancelNFSeWithReasonHandler - POST /api/requests/{id}/nfse/cancelar
func CancelNFSeWithReasonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	requestIdStr := r.PathValue("id")

	var req struct {
		Motivo        int    `json:"motivo"`
		Justificativa string `json:"justificativa"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Update NFSe status
	if err := shared.GetDB().Model(&shared.NFSe{}).
		Where("request_id = ?", requestIdStr).
		Updates(map[string]interface{}{
			"status":        "cancelado",
			"motivo_cancel": req.Justificativa,
		}).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to cancel NFSe")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Solicitação de cancelamento enviada"})
}

// GetDANFSeHandler - GET /api/requests/{id}/nfse/danfse
func GetDANFSeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	// Return a sample PDF URL or redirect
	http.Redirect(w, r, "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf", http.StatusFound)
}

// GetEventosHandler - GET /api/requests/{id}/nfse/eventos
func GetEventosHandler(w http.ResponseWriter, r *http.Request) {
	// Stub return empty list
	shared.SuccessResponse(w, []string{})
}
