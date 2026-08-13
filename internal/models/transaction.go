package models

import "time"

// TransactionType merepresentasikan 4 tipe transaksi sesuai Bab 21 RPS.
type TransactionType string

const (
	TransactionSale        TransactionType = "SALE"
	TransactionPurchase    TransactionType = "PURCHASE"
	TransactionOperational TransactionType = "OPERATIONAL"
	TransactionExpense     TransactionType = "EXPENSE"
)

// Transaction merepresentasikan tabel `transactions`.
type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	Type            TransactionType   `gorm:"size:20;not null" json:"type"`
	TransactionDate time.Time         `gorm:"not null" json:"transaction_date"`
	Description     *string           `gorm:"size:255" json:"description"`
	TotalAmount     float64           `gorm:"not null" json:"total_amount"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
}

// TransactionItem merepresentasikan tabel `transaction_items`.
// Hanya dipakai untuk tipe SALE dan PURCHASE (lihat hasItems di Flutter).
type TransactionItem struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `json:"-"`
	ProductID     uint    `json:"product_id"`
	ProductName   string  `gorm:"size:150;not null" json:"product_name"`
	Quantity      int     `gorm:"not null" json:"quantity"`
	UnitPrice     float64 `gorm:"not null" json:"unit_price"`
}