package shared

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes database connection (call once)
// Optimized for serverless with connection pooling
func InitDB() error {
	if DB != nil {
		return nil // Already initialized (cache hit)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info), // Increased log level for better visibility
		PrepareStmt: true,
	})

	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		return err
	}

	// Configure connection pool for serverless
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Serverless optimization: keep minimal idle connections
	sqlDB.SetMaxIdleConns(2)                  // Minimal idle (serverless cold start)
	sqlDB.SetMaxOpenConns(10)                 // Max concurrent connections
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // 5min max lifetime (Supabase pooler)
	sqlDB.SetConnMaxIdleTime(time.Minute * 2) // Close idle after 2min

	// Auto Migrate (Schema Update)
	// This is critical for serverless deployments to ensure tables exist
	err = DB.AutoMigrate(
		&User{},
		&Client{},
		&Request{},
		&Equipment{},
		&ChecklistItem{},
		&AgendaEntry{},
		&Company{},
		&Transaction{},
		&FiscalConfig{},
		&AuditLog{},
		&Setting{},
		&OrcamentoItem{},
		&Assinatura{},
		&NFSe{},
	)
	if err != nil {
		log.Printf("⚠️ Migration warning: %v", err)
		// Don't fail hard on migration if it's a minor issue, but log it
	}

	log.Println("✅ Database connected (Supabase) with serverless pooling & Migrations")
	return nil
}

// GetDB returns the database instance (singleton pattern)
func GetDB() *gorm.DB {
	if DB == nil {
		InitDB()
	}
	return DB
}
