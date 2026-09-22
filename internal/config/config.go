package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua environment variable yang dibutuhkan aplikasi.
type Config struct {
	DatabaseURL string
	AppPort     string
}

// Load membaca file .env dan mengembalikan Config yang siap dipakai.
func Load() *Config {
	// Kalau .env tidak ditemukan, tidak fatal — bisa jadi environment
	// variable sudah di-set langsung di sistem (misal saat deploy nanti).
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: file .env tidak ditemukan, menggunakan environment variable sistem")
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AppPort:     os.Getenv("APP_PORT"),
	}
}
