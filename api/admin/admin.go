package admin

import (
	"inovar/lib/shared"
	"net/http"
)

// WipeHandler triggers a full database wipe and re-seed.
// WARNING: This is dangerous.
func WipeHandler(w http.ResponseWriter, r *http.Request) {
	// Simple protection: Check for a secret header or query param
	// Since asking user to set header is hard in browser, use query param ?secret=force_wipe_2026
	secret := r.URL.Query().Get("secret")
	if secret != "force_wipe_2026" {
		shared.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Run SeedAdmin synchronously to report errors
	// (We might want to refactor SeedAdmin to return error)
	// For now, SeedAdmin just logs. We trust it works or logs.

	// We can't easily capture SeedAdmin logs unless we change the function signature.
	// I'll update SeedAdmin signature in a bit?
	// Or just call it and assume success if no panic.

	// Better: Copy logic or Refactor SeedAdmin to return error.
	// I'll refactor SeedAdmin to return error.

	if err := shared.SeedAdmin(); err != nil {
		shared.ErrorResponse(w, http.StatusInternalServerError, "Wipe failed: "+err.Error())
		return
	}

	shared.SuccessResponse(w, map[string]string{
		"message":  "Database wiped and admin seeded successfully",
		"admin":    "admin@inovar.com",
		"password": "123456",
	})
}
