package repositories

import (
	"time"

	"gorm.io/gorm"

	"sellora_backend/internal/models"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// SummaryResult menampung angka mentah hasil kalkulasi, sebelum
// dikonversi ke DTO response oleh handler.
type SummaryResult struct {
	Revenue            float64
	COGS               float64
	GrossProfit        float64
	OperationalExpense float64
	OtherExpense       float64
	NetProfit          float64
	ProfitMargin       float64
}

// GetSummary menghitung ringkasan bisnis untuk rentang tanggal tertentu
// (Bab 15 RPS - Rumus Bisnis).
func (r *ReportRepository) GetSummary(start, end time.Time) (SummaryResult, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Items").
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Find(&transactions).Error
	if err != nil {
		return SummaryResult{}, err
	}

	var revenue, cogs, operationalExpense, otherExpense float64

	for _, t := range transactions {
		switch t.Type {
		case models.TransactionSale:
			revenue += t.TotalAmount
			// HPP dihitung dari cost_price produk SAAT INI, bukan
			// dari data historis -- ini simplifikasi yang cukup untuk
			// skala aplikasi personal seperti Sellora.
			for _, item := range t.Items {
				var product models.Product
				if err := r.db.First(&product, item.ProductID).Error; err == nil {
					cogs += product.CostPrice * float64(item.Quantity)
				}
			}
		case models.TransactionOperational:
			operationalExpense += t.TotalAmount
		case models.TransactionExpense:
			otherExpense += t.TotalAmount
		}
	}

	grossProfit := revenue - cogs
	netProfit := grossProfit - operationalExpense - otherExpense
	margin := 0.0
	if revenue > 0 {
		margin = (netProfit / revenue) * 100
	}

	return SummaryResult{
		Revenue:            revenue,
		COGS:               cogs,
		GrossProfit:        grossProfit,
		OperationalExpense: operationalExpense,
		OtherExpense:       otherExpense,
		NetProfit:          netProfit,
		ProfitMargin:       margin,
	}, nil
}

// ChartPointResult untuk grafik tren per hari.
type ChartPointResult struct {
	Date    time.Time
	Revenue float64
	Profit  float64
}

// GetDailyTrend menghitung Pendapatan & Laba per hari, dalam rentang
// tanggal tertentu. Dipakai untuk grafik Dashboard & Report.
func (r *ReportRepository) GetDailyTrend(start, end time.Time) ([]ChartPointResult, error) {
	// Buat map tanggal -> summary, supaya hari tanpa transaksi tetap
	// muncul di grafik dengan nilai 0 (bukan hilang dari daftar).
	dayCount := int(end.Sub(start).Hours()/24) + 1
	results := make([]ChartPointResult, dayCount)
	for i := 0; i < dayCount; i++ {
		results[i] = ChartPointResult{Date: start.AddDate(0, 0, i)}
	}

	var transactions []models.Transaction
	err := r.db.Preload("Items").
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	for _, t := range transactions {
		dayIndex := int(t.TransactionDate.Sub(start).Hours() / 24)
		if dayIndex < 0 || dayIndex >= len(results) {
			continue
		}

		if t.Type == models.TransactionSale {
			results[dayIndex].Revenue += t.TotalAmount

			var cogs float64
			for _, item := range t.Items {
				var product models.Product
				if err := r.db.First(&product, item.ProductID).Error; err == nil {
					cogs += product.CostPrice * float64(item.Quantity)
				}
			}
			results[dayIndex].Profit += t.TotalAmount - cogs
		}
	}

	return results, nil
}

// BreakdownResult untuk breakdown pengeluaran per kategori/tipe.
type BreakdownResult struct {
	Label  string
	Amount float64
}

// GetExpenseBreakdown mengelompokkan pengeluaran: PURCHASE per kategori
// produk, OPERATIONAL & EXPENSE sebagai baris tersendiri.
func (r *ReportRepository) GetExpenseBreakdown(start, end time.Time) ([]BreakdownResult, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Items").
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	categoryTotals := make(map[string]float64)
	var operationalTotal, expenseTotal float64

	for _, t := range transactions {
		switch t.Type {
		case models.TransactionPurchase:
			for _, item := range t.Items {
				var product models.Product
				if err := r.db.Preload("Category").First(&product, item.ProductID).Error; err == nil {
					categoryTotals[product.Category.Name] += item.UnitPrice * float64(item.Quantity)
				}
			}
		case models.TransactionOperational:
			operationalTotal += t.TotalAmount
		case models.TransactionExpense:
			expenseTotal += t.TotalAmount
		}
	}

	var results []BreakdownResult
	for label, amount := range categoryTotals {
		results = append(results, BreakdownResult{Label: label, Amount: amount})
	}
	if operationalTotal > 0 {
		results = append(results, BreakdownResult{Label: "Operasional", Amount: operationalTotal})
	}
	if expenseTotal > 0 {
		results = append(results, BreakdownResult{Label: "Pengeluaran Lain", Amount: expenseTotal})
	}

	return results, nil
}