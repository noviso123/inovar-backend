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
	// Mock routes list or dynamic extraction
	shared.SuccessResponse(w, []string{
		"/api/users", "/api/clients", "/api/requests", "/api/equipments",
		"/api/agenda", "/api/company", "/api/finance", "/api/audit",
		"/api/settings", "/api/fiscal", "/api/system",
	})
}

func TablesHandler(w http.ResponseWriter, r *http.Request) {
	// Return the list of Portuguese tables from our models (or DB directly)
	tables := []string{
		"users", "clientes", "solicitacoes", "equipamentos",
		"checklists", "agenda", "empresas", "financeiro",
		"audit_logs", "configuracoes", "orcamentos", "assinaturas", "nfse",
	}
	shared.SuccessResponse(w, tables)
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
