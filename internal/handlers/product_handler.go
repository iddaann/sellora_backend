package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sellora_backend/internal/handlers/dto"
	"sellora_backend/internal/models"
	"sellora_backend/internal/repositories"
)

type ProductHandler struct {
	repo *repositories.ProductRepository
}

func NewProductHandler(repo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GetAll -> GET /api/products
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}
	c.JSON(http.StatusOK, dto.FromModelList(products))
}

// createProductRequest merepresentasikan body request untuk Create & Update.
type createProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
	SellPrice  float64 `json:"sell_price" binding:"required,gt=0"`
	CostPrice  float64 `json:"cost_price" binding:"required,gte=0"`
	Stock      int     `json:"stock"`
	Unit       string  `json:"unit" binding:"required"`
}

// Create -> POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		SellPrice:  req.SellPrice,
		CostPrice:  req.CostPrice,
		Stock:      req.Stock,
		Unit:       req.Unit,
	}

	if err := h.repo.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}

	// Ambil ulang dengan Category ter-preload supaya response lengkap.
	created, err := h.repo.GetByID(product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Produk tersimpan tapi gagal dimuat ulang"})
		return
	}
	c.JSON(http.StatusCreated, dto.FromModel(created))
}

// Update -> PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		ID:         uint(id),
		Name:       req.Name,
		CategoryID: req.CategoryID,
		SellPrice:  req.SellPrice,
		CostPrice:  req.CostPrice,
		Stock:      req.Stock,
		Unit:       req.Unit,
	}

	if err := h.repo.Update(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui produk"})
		return
	}

	updated, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Produk tersimpan tapi gagal dimuat ulang"})
		return
	}
	c.JSON(http.StatusOK, dto.FromModel(updated))
}

// Delete -> DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}