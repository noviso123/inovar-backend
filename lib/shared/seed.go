package shared

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedAdmin checks if admin user exists, if not creates it.
func SeedAdmin() {
	// Ensure DB connection is available
	db := GetDB()

	email := "admin@inovar.com"
	password := "123456"

	var user User
	// Check if user exists
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		log.Printf("✅ Admin user already exists: %s", email)
		return
	}

	// Create user
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	active := true
	mustChange := false

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
		return "00000000-0000-4000-8000-000000000001" // Fallback if random fails (extremely unlikely)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	return hex.EncodeToString(b[:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}
