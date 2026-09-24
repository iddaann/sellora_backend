package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"sellora_backend/internal/handlers/dto"
	"sellora_backend/internal/repositories"
)

type ReportHandler struct {
	repo *repositories.ReportRepository
}

func NewReportHandler(repo *repositories.ReportRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

// resolveDateRange menerjemahkan query param "period" (daily/weekly/monthly)
// menjadi rentang tanggal, sesuai ReportPeriod.range di Flutter.
func resolveDateRange(period string) (time.Time, time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch period {
	case "weekly":
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := today.AddDate(0, 0, -(weekday - 1))
		end := start.AddDate(0, 0, 6)
		return start, end.Add(24*time.Hour - time.Second)
	case "monthly":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return start, end.Add(24*time.Hour - time.Second)
	default:
		return today, today.Add(24*time.Hour - time.Second)
	}
}

func periodLabel(period string) string {
	switch period {
	case "weekly":
		return "Mingguan"
	case "monthly":
		return "Bulanan"
	default:
		return "Harian"
	}
}

// GetSummary -> GET /api/dashboard atau GET /api/reports/summary?period=weekly
func (h *ReportHandler) GetSummary(c *gin.Context) {
	period := c.DefaultQuery("period", "daily")
	start, end := resolveDateRange(period)

	result, err := h.repo.GetSummary(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung ringkasan"})
		return
	}

	c.JSON(http.StatusOK, dto.SummaryResponse{
		PeriodLabel:        periodLabel(period),
		Revenue:            result.Revenue,
		COGS:               result.COGS,
		GrossProfit:        result.GrossProfit,
		OperationalExpense: result.OperationalExpense,
		OtherExpense:       result.OtherExpense,
		NetProfit:          result.NetProfit,
		ProfitMargin:       result.ProfitMargin,
	})
}

// GetTrendChart -> GET /api/reports/chart?period=weekly
func (h *ReportHandler) GetTrendChart(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	start, end := resolveDateRange(period)

	points, err := h.repo.GetDailyTrend(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung grafik"})
		return
	}

	response := make([]dto.ChartPointResponse, len(points))
	for i, p := range points {
		response[i] = dto.ChartPointResponse{
			Date:    p.Date.Format("2006-01-02"),
			Revenue: p.Revenue,
			Profit:  p.Profit,
		}
	}
	c.JSON(http.StatusOK, response)
}

// GetExpenseBreakdown -> GET /api/reports/breakdown?period=weekly
func (h *ReportHandler) GetExpenseBreakdown(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	start, end := resolveDateRange(period)

	items, err := h.repo.GetExpenseBreakdown(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung breakdown"})
		return
	}

	var total float64
	for _, item := range items {
		total += item.Amount
	}

	response := make([]dto.ExpenseBreakdownResponse, len(items))
	for i, item := range items {
		percentage := 0.0
		if total > 0 {
			percentage = (item.Amount / total) * 100
		}
		response[i] = dto.ExpenseBreakdownResponse{
			Label:      item.Label,
			Amount:     item.Amount,
			Percentage: percentage,
		}
	}
	c.JSON(http.StatusOK, response)
}
