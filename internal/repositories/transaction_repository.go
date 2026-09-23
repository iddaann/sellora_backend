package repositories

import (
	"errors"

	"gorm.io/gorm"

	"sellora_backend/internal/models"
)

var ErrInsufficientStock = errors.New("stok tidak mencukupi")
var ErrProductNotFound = errors.New("produk tidak ditemukan")

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetAll(filterType string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Preload("Items").Order("transaction_date DESC")

	if filterType != "" {
		query = query.Where("type = ?", filterType)
	}

	err := query.Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i := range transaction.Items {
			var product models.Product
			if err := tx.First(&product, transaction.Items[i].ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrProductNotFound
				}
				return err
			}
			transaction.Items[i].ProductName = product.Name
			transaction.Items[i].CostPrice = product.CostPrice
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		if transaction.Type != models.TransactionSale &&
			transaction.Type != models.TransactionPurchase {
			return nil
		}

		for _, item := range transaction.Items {
			if transaction.Type == models.TransactionSale {
				result := tx.Model(&models.Product{}).
					Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
					UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))

				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					var product models.Product
					if err := tx.First(&product, item.ProductID).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return ErrProductNotFound
						}
						return err
					}
					return ErrInsufficientStock
				}
				continue
			}

			result := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrProductNotFound
			}
		}

		return nil
	})
}
