package repositories

import (
	"gorm.io/gorm"

	"sellora_backend/internal/models"
)

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
		query = query.Where("type= ?", filterType)
	}

	err := query.Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		for _, item := range transaction.Items {
			delta := item.Quantity
			if transaction.Type == models.TransactionSale {
				delta = -delta
			}

			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock = ?", delta)).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
}