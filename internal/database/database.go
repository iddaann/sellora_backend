package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sellora_backend/internal/config"
)

func Connect(cfg *config.Config) *gorm.DB {
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diisi di file .env")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("gagal connect ke database Supabase: %v", err)
	}

	log.Println("berhasil connect ke database Supabase PostgreSQL")
	return db
}
