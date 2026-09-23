package handlers

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"sellora_backend/internal/handlers/dto"
	"sellora_backend/internal/models"
	"sellora_backend/internal/repositories"
)

type TransactionHandler struct {
	repo *repositories.TransactionRepository
}

func NewTransactionHandler(repo *repositories.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

// GetAll -> GET /api/transactions?type=SALE (query param "type" opsional)
func (h *TransactionHandler) GetAll(c *gin.Context) {
	filterType := c.Query("type")

	transactions, err := h.repo.GetAll(filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}
	c.JSON(http.StatusOK, dto.FromTransactionModelList(transactions))
}

type createTransactionItemRequest struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity" binding:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
}

type createTransactionRequest struct {
	Type            string                          `json:"type" binding:"required,oneof=SALE PURCHASE OPERATIONAL EXPENSE"`
	TransactionDate string                          `json:"transaction_date" binding:"required"`
	Description     *string                         `json:"description"`
	TotalAmount     float64                          `json:"total_amount" binding:"required,gt=0"`
	Items           []createTransactionItemRequest  `json:"items"`
}

// Create -> POST /api/transactions
func (h *TransactionHandler) Create(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse tanggal dari format ISO8601 yang dikirim Flutter
	// (DateTime.toIso8601String() di transaction.dart Step 5.3).
	parsedDate, err := parseISODate(req.TransactionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format transaction_date tidak valid"})
		return
	}

	items := make([]models.TransactionItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = models.TransactionItem{
			ProductID:   itemReq.ProductID,
			ProductName: itemReq.ProductName,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
		}
	}

	// SALE dan PURCHASE wajib memiliki item karena item digunakan untuk
	// menghitung total dan memperbarui stok.
	if (req.Type == "SALE" || req.Type == "PURCHASE") && len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi SALE/PURCHASE harus memiliki minimal satu item"})
		return
	}

	// OPERATIONAL dan EXPENSE tidak berhubungan dengan produk.
	if (req.Type == "OPERATIONAL" || req.Type == "EXPENSE") && len(req.Items) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi OPERATIONAL/EXPENSE tidak boleh memiliki item produk"})
		return
	}

	// Untuk SALE/PURCHASE, backend menghitung ulang total dari detail transaksi.
	if req.Type == "SALE" || req.Type == "PURCHASE" {
		calculatedTotal := 0.0
		for _, item := range req.Items {
			calculatedTotal += float64(item.Quantity) * item.UnitPrice
		}

		if math.Abs(calculatedTotal-req.TotalAmount) > 0.01 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":            "total_amount tidak sesuai dengan detail transaksi",
				"calculated_total": calculatedTotal,
			})
			return
		}
	}

	transaction := models.Transaction{
		Type:            models.TransactionType(req.Type),
		TransactionDate: parsedDate,
		Description:     req.Description,
		TotalAmount:     req.TotalAmount,
		Items:           items,
	}

	if err := h.repo.Create(&transaction); err != nil {
		if errors.Is(err, repositories.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak mencukupi"})
			return
		}
		if errors.Is(err, repositories.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusCreated, dto.FromTransactionModel(transaction))
}