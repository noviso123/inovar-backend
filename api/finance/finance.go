package finance

import (
	"encoding/json"
	"inovar/lib/shared"
	"net/http"
	"time"
)

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
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

	// Calculate summary logic
	// For now, let's mock or calculate basic sums
	var income float64
	var expense float64

	shared.GetDB().Model(&shared.Transaction{}).Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Scan(&income)
	shared.GetDB().Model(&shared.Transaction{}).Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0)").Scan(&expense)

	balance := income - expense

	// Get recent transactions
	var transactions []shared.Transaction
	shared.GetDB().Order("date DESC").Limit(10).Find(&transactions)

	response := map[string]interface{}{
		"income":       income,
		"expense":      expense,
		"balance":      balance,
		"transactions": transactions,
	}

	shared.SuccessResponse(w, response)
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
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

	var transactions []shared.Transaction
	// check date filters
	if err := shared.GetDB().Order("date DESC").Limit(100).Find(&transactions).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Query failed")
		return
	}

	shared.SuccessResponse(w, transactions)
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// ... (Standard Create Logic)
	// Implement if needed by frontend (frontend didn't show createTransaction explicitly in apiService but implied)
	// Let's implement basics

	// Decoding...
	var t shared.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		shared.ErrorResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}
	if t.Date.IsZero() {
		t.Date = time.Now()
	}

	if err := shared.GetDB().Create(&t).Error; err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Create failed")
		return
	}
	shared.SuccessResponse(w, t)
}
