package budget

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
	"strconv"
)

// ListItemsHandler - GET /api/requests/{requestId}/orcamento/itens (Not explicitly in apiService but good to have)
// Actually apiService uses getRequest which likely includes items?
// apiService has getOrcamentoSugestoes and addOrcamentoItem.
// Let's implement what's needed.

// SugestoesHandler - GET /api/requests/orcamento/sugestoes
func SugestoesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Mock suggestions for now - in a real app this would query common items or a catalog
	sugestoes := []map[string]interface{}{
		{"descricao": "Visita Técnica", "valorUnit": 150.00, "tipo": "mao_de_obra"},
		{"descricao": "Cabo de Rede (Metro)", "valorUnit": 5.00, "tipo": "peca"},
		{"descricao": "Conector RJ45", "valorUnit": 2.50, "tipo": "peca"},
		{"descricao": "Instalação de Câmera", "valorUnit": 120.00, "tipo": "mao_de_obra"},
		{"descricao": "Configuração de DVR", "valorUnit": 200.00, "tipo": "mao_de_obra"},
	}

	shared.SuccessResponse(w, sugestoes)
}

// AddItemHandler - POST /api/requests/{requestId}/orcamento/itens
func AddItemHandler(w http.ResponseWriter, r *http.Request) {
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

	var item shared.OrcamentoItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item.RequestID = uint(requestID)

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Create(&item).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to create budget item")
		return
	}

	shared.SuccessResponse(w, item)
}

// RemoveItemHandler - DELETE /api/requests/{requestId}/orcamento/itens/{itemId}
func RemoveItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	itemIdStr := r.PathValue("itemId")

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := shared.GetDB().Delete(&shared.OrcamentoItem{}, itemIdStr).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AprovarHandler - POST /api/requests/{requestId}/orcamento/aprovar
func AprovarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		shared.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	requestIdStr := r.PathValue("id")

	if err := shared.InitDB(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Mark all items for this request as approved
	// In a more complex system, maybe we approve the whole budget state on the Request model
	// For now, let's update the items logic if that's what's intended, OR update the Request status

	// Let's assume hitting this endpoint approves the budget and allows work to proceed.
	// We might want to update Request status to 'in_progress' or similar?
	// Or just mark items? The model has 'Aprovado' bool.

	result := shared.GetDB().Model(&shared.OrcamentoItem{}).
		Where("request_id = ?", requestIdStr).
		Update("aprovado", true)

	if result.Error != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Failed to approve budget")
		return
	}

	shared.SuccessResponse(w, map[string]string{"message": "Orçamento aprovado com sucesso"})
}
