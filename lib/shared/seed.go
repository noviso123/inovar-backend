package shared

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedAdmin checks if admin user exists, if not creates it.
// If it exists, it UPDATES the password to ensure access.
func SeedAdmin() {
	// Ensure DB connection is available
	db := GetDB()

	// ⚠️ DANGER: Wipe all data as requested by user (Full Reset)
	log.Println("🔥 STARTING FULL DATABASE WIPE (TRUNCATE CASCADE)...")
	// Execute dynamic TRUNCATE on all tables in public schema
	if err := db.Exec(`
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`).Error; err != nil {
		log.Printf("❌ Failed to wipe database: %v", err)
	} else {
		log.Println("✅ Database wiped successfully (All tables truncated)")
	}

	email := "admin@inovar.com"
	password := "123456"

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	active := true
	mustChange := false

	var user User
	// Check if user exists (should not exist after wipe, but safe coding)
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		// User exists, FORCE UPDATE password and active status
		if err := db.Model(&user).Updates(map[string]interface{}{
			"password_hash":        string(hashed),
			"active":               true,
			"must_change_password": false,
			"role":                 "ADMIN_SISTEMA", // Maintain role
		}).Error; err != nil {
			log.Printf("❌ Failed to update admin user: %v", err)
		} else {
			log.Printf("✅ Admin user credentials updated (password reset): %s", email)
		}
		return
	}

	// Create user if not exists
	newUser := User{
		ID:                 generateSeedUUID(),
		Name:               "Admin Inovar",
		Email:              email,
		PasswordHash:       string(hashed),
		Role:               "ADMIN_SISTEMA",
		Active:             &active,
		MustChangePassword: &mustChange,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Printf("❌ Failed to seed admin user: %v", err)
	} else {
		log.Printf("🚀 Admin user seeded successfully: %s", email)
	}
}

func generateSeedUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000-0000-4000-8000-000000000001" // Fallback
	}
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	return hex.EncodeToString(b[:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}
