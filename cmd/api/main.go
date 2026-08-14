package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"sellora_backend/internal/config"
	"sellora_backend/internal/database"
	"sellora_backend/internal/handlers"
	"sellora_backend/internal/models"
	"sellora_backend/internal/repositories"
	"sellora_backend/internal/routes"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal ambil koneksi SQL: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping database gagal: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionItem{},
	); err != nil {
		log.Fatalf("Migrasi database gagal: %v", err)
	}
	log.Println("Migrasi database berhasil")

	// --- Rakit dependency: Repository -> Handler ---
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	productRepo := repositories.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionHandler := handlers.NewTransactionHandler(transactionRepo)

	reportRepo := repositories.NewReportRepository(db)
	reportHandler := handlers.NewReportHandler(reportRepo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Sellora API is running",
		})
	})

	routes.Setup(r, categoryHandler, productHandler, transactionHandler, reportHandler)

	log.Printf("Server jalan di port %s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}