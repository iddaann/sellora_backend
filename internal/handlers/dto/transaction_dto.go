package dto

import (
	"time"

	"sellora_backend/internal/models"
)

type TransactionItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type TransactionResponse struct {
	ID              uint                       `json:"id"`
	Type            string                     `json:"type"`
	TransactionDate time.Time                  `json:"transaction_date"`
	Description     *string                    `json:"description"`
	TotalAmount     float64                    `json:"total_amount"`
	Items           []TransactionItemResponse  `json:"items"`
}

func FromTransactionModel(t models.Transaction) TransactionResponse {
	items := make([]TransactionItemResponse, len(t.Items))
	for i, item := range t.Items {
		items[i] = TransactionItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	return TransactionResponse{
		ID:              t.ID,
		Type:            string(t.Type),
		TransactionDate: t.TransactionDate,
		Description:     t.Description,
		TotalAmount:     t.TotalAmount,
		Items:           items,
	}
}

func FromTransactionModelList(transactions []models.Transaction) []TransactionResponse {
	result := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		result[i] = FromTransactionModel(t)
	}
	return result
}