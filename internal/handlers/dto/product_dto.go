package dto

import "sellora_backend/internal/models"

// ProductResponse adalah bentuk JSON yang dikirim ke client.
// Field "category" sengaja string (nama kategori), bukan objek Category
// penuh, supaya cocok dengan model Flutter Product.fromJson yang
// mengharapkan json['category'] as String.
type ProductResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	CategoryID uint    `json:"category_id"`
	SellPrice  float64 `json:"sell_price"`
	CostPrice  float64 `json:"cost_price"`
	Stock      int     `json:"stock"`
	Unit       string  `json:"unit"`
}

// FromModel mengonversi models.Product (hasil query database, lengkap
// dengan relasi Category yang sudah di-preload) menjadi ProductResponse.
func FromModel(p models.Product) ProductResponse {
	return ProductResponse{
		ID:         p.ID,
		Name:       p.Name,
		Category:   p.Category.Name,
		CategoryID: p.CategoryID,
		SellPrice:  p.SellPrice,
		CostPrice:  p.CostPrice,
		Stock:      p.Stock,
		Unit:       p.Unit,
	}
}

// FromModelList mengonversi banyak Product sekaligus.
func FromModelList(products []models.Product) []ProductResponse {
	result := make([]ProductResponse, len(products))
	for i, p := range products {
		result[i] = FromModel(p)
	}
	return result
}