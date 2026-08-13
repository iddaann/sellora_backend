package models

type Product struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `gorm:"size:150;not null" json:"name"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"-"`
	SellPrice  float64  `gorm:"not null" json:"sell_price"`
	CostPrice  float64  `gorm:"not null" json:"cost_price"`
	Stock      int      `gorm:"not null" json:"stock"`
	Unit       string   `gorm:"size:20;not null" json:"unit"`
}

func (p Product) MarshallCategoryName() string {
	return p.Category.Name
}