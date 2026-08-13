package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sellora_backend/internal/models"
	"sellora_backend/internal/repositories"
)

// CategoryHandler menerima HTTP request dan menjembatani ke repository.
type CategoryHandler struct {
	repo *repositories.CategoryRepository
}

func NewCategoryHandler(repo *repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GetAll -> GET /api/categories
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// Create -> POST /api/categories
type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori wajib diisi"})
		return
	}

	category := models.Category{Name: req.Name}
	if err := h.repo.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// Delete -> DELETE /api/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}