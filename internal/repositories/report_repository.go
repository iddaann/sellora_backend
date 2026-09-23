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

type SummaryResult struct {
	Revenue            float64
	COGS               float64
	GrossProfit        float64
	OperationalExpense float64
	OtherExpense       float64
	NetProfit          float64
	ProfitMargin       float64
}

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
			for _, item := range t.Items {
				cogs += item.CostPrice * float64(item.Quantity)
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

type ChartPointResult struct {
	Date    time.Time
	Revenue float64
	Profit  float64
}

func (r *ReportRepository) GetDailyTrend(start, end time.Time) ([]ChartPointResult, error) {
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
				cogs += item.CostPrice * float64(item.Quantity)
			}
			results[dayIndex].Profit += t.TotalAmount - cogs
		}
	}

	return results, nil
}

type BreakdownResult struct {
	Label  string
	Amount float64
}

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
