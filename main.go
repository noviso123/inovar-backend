package main

import (
	"log"
	"net/http"
	"os"

	"inovar/api/agenda"
	"inovar/api/audit"
	"inovar/api/budget"
	"inovar/api/checklists"
	"inovar/api/clients"
	"inovar/api/company"
	"inovar/api/equipments"
	"inovar/api/finance"
	"inovar/api/fiscal"
	"inovar/api/login"
	"inovar/api/requests"
	"inovar/api/settings"
	"inovar/api/signatures"
	"inovar/api/system"
	"inovar/api/upload"
	"inovar/api/users"
	"inovar/lib/shared"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables if .env exists (local dev)
	_ = godotenv.Load()

	// Initialize Database
	log.Println("Initializing database...")
	if err := shared.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Define Routes (Go 1.22+)
	mux := http.NewServeMux()

	// Catch-all / Health Check (Returns JSON error for invalid routes to avoid "Unexpected token M")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == "GET" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Inovar Backend Running"))
			return
		}
		shared.ErrorResponse(w, http.StatusNotFound, "Route not found or method not allowed")
	})

	// Auth
	mux.HandleFunc("POST /api/login", login.LoginHandler)
	mux.HandleFunc("GET /api/auth/me", login.MeHandler)
	mux.HandleFunc("PUT /api/auth/me", login.UpdateProfileHandler)
	mux.HandleFunc("POST /api/auth/change-password", login.ChangePasswordHandler)
	mux.HandleFunc("POST /api/auth/forgot-password", login.ForgotPasswordHandler)
	mux.HandleFunc("POST /api/auth/reset-password", login.ResetPasswordHandler)

	// Clients
	mux.HandleFunc("GET /api/clients", clients.ListHandler)
	mux.HandleFunc("POST /api/clients", clients.CreateHandler)
	mux.HandleFunc("GET /api/clients/{id}", clients.GetHandler)
	mux.HandleFunc("PUT /api/clients/{id}", clients.UpdateHandler)
	mux.HandleFunc("DELETE /api/clients/{id}", clients.DeleteHandler)

	// Users
	mux.HandleFunc("GET /api/users", users.ListHandler)
	mux.HandleFunc("POST /api/users", users.CreateHandler)
	mux.HandleFunc("GET /api/users/{id}", users.GetHandler)
	mux.HandleFunc("PUT /api/users/{id}", users.UpdateHandler)
	mux.HandleFunc("PATCH /api/users/{id}/block", users.BlockUserHandler)
	mux.HandleFunc("DELETE /api/users/{id}", users.DeleteHandler)
	mux.HandleFunc("POST /api/users/{id}/reset-password", users.AdminResetPasswordHandler)

	// Requests
	mux.HandleFunc("GET /api/requests", requests.ListHandler)
	mux.HandleFunc("POST /api/requests", requests.CreateHandler)
	mux.HandleFunc("GET /api/requests/{id}", requests.GetHandler)
	mux.HandleFunc("PATCH /api/requests/{id}/status", requests.UpdateStatusHandler)
	mux.HandleFunc("DELETE /api/requests/{id}", requests.DeleteHandler)
	mux.HandleFunc("POST /api/requests/{id}/attachments", requests.UploadAttachmentHandler)
	mux.HandleFunc("DELETE /api/requests/{id}/attachments/{fileId}", requests.DeleteAttachmentHandler)

	// Upload
	mux.HandleFunc("POST /api/upload", upload.Handler)

	// Equipments
	mux.HandleFunc("GET /api/equipments", equipments.ListHandler)
	mux.HandleFunc("POST /api/equipments", equipments.CreateHandler)
	mux.HandleFunc("GET /api/equipments/{id}", equipments.GetHandler)
	mux.HandleFunc("PUT /api/equipments/{id}", equipments.UpdateHandler)
	mux.HandleFunc("PATCH /api/equipments/{id}/deactivate", equipments.DeactivateHandler)
	mux.HandleFunc("PATCH /api/equipments/{id}/reactivate", equipments.ReactivateHandler)
	mux.HandleFunc("DELETE /api/equipments/{id}", equipments.DeleteHandler)

	// Checklists
	mux.HandleFunc("GET /api/requests/{requestId}/checklists", checklists.ListHandler)
	mux.HandleFunc("POST /api/requests/{requestId}/checklists", checklists.CreateHandler)
	mux.HandleFunc("PATCH /api/requests/{requestId}/checklists/{itemId}", checklists.UpdateHandler)
	mux.HandleFunc("DELETE /api/requests/{requestId}/checklists/{itemId}", checklists.DeleteHandler)

	// Agenda
	mux.HandleFunc("GET /api/agenda", agenda.ListHandler)
	mux.HandleFunc("POST /api/agenda", agenda.CreateHandler)
	mux.HandleFunc("PUT /api/agenda/{id}", agenda.UpdateHandler)
	mux.HandleFunc("DELETE /api/agenda/{id}", agenda.DeleteHandler)

	// Company
	mux.HandleFunc("GET /api/company", company.GetHandler)
	mux.HandleFunc("PUT /api/company", company.UpdateHandler)

	// Finance
	mux.HandleFunc("GET /api/finance/summary", finance.SummaryHandler)
	mux.HandleFunc("GET /api/finance/transactions", finance.TransactionsHandler)
	mux.HandleFunc("POST /api/finance/transactions", finance.CreateTransactionHandler)

	// Audit
	mux.HandleFunc("GET /api/audit", audit.ListHandler)

	// Settings
	mux.HandleFunc("GET /api/settings", settings.GetHandler)
	mux.HandleFunc("PUT /api/settings", settings.UpdateHandler)

	// Fiscal & NFSe
	mux.HandleFunc("GET /api/fiscal/config", fiscal.ConfigHandler)
	mux.HandleFunc("POST /api/fiscal/certificate", fiscal.CertificateHandler)
	mux.HandleFunc("POST /api/requests/{id}/nfse", fiscal.IssueNFSeHandler)
	mux.HandleFunc("GET /api/requests/{id}/nfse", fiscal.GetNFSeHandler)
	mux.HandleFunc("DELETE /api/requests/{id}/nfse", fiscal.CancelNFSeHandler)
	mux.HandleFunc("POST /api/requests/{id}/nfse/cancelar", fiscal.CancelNFSeWithReasonHandler)
	mux.HandleFunc("GET /api/requests/{id}/nfse/danfse", fiscal.GetDANFSeHandler)
	mux.HandleFunc("GET /api/requests/{id}/nfse/eventos", fiscal.GetEventosHandler)
	mux.HandleFunc("GET /api/fiscal/regimes", fiscal.RegimesHandler)
	mux.HandleFunc("GET /api/fiscal/lookup/{cnpj}", fiscal.LookupCNPJHandler)

	// Budget (Orcamento)
	mux.HandleFunc("GET /api/requests/orcamento/sugestoes", budget.SugestoesHandler)
	mux.HandleFunc("POST /api/requests/{id}/orcamento/itens", budget.AddItemHandler)
	mux.HandleFunc("DELETE /api/requests/{id}/orcamento/itens/{itemId}", budget.RemoveItemHandler)
	mux.HandleFunc("POST /api/requests/{id}/orcamento/aprovar", budget.AprovarHandler)

	// Signatures
	mux.HandleFunc("POST /api/requests/{id}/assinatura", signatures.SaveHandler)

	// System
	mux.HandleFunc("GET /api/system/whatsapp", system.WhatsAppStatusHandler)
	mux.HandleFunc("GET /api/system/routes", system.RoutesHandler)

	// Wrap with CORS and Logging
	handler := corsMiddleware(loggingMiddleware(mux))

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Allow specific origins including Vercel and Localhost
		allowedOrigins := map[string]bool{
			"https://inovar-gestao.vercel.app": true,
			"http://localhost:5173":            true,
			"http://localhost:3000":            true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Fallback or development
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
