package dto

// SummaryResponse dipakai untuk endpoint Dashboard maupun Report --
// strukturnya sama persis dengan DashboardSummary & ReportSummary
// di Flutter (Step 4.1 & 7.1).
type SummaryResponse struct {
	Revenue            float64 `json:"revenue"`
	COGS               float64 `json:"cogs"`
	GrossProfit        float64 `json:"gross_profit"`
	OperationalExpense float64 `json:"operational_expense"`
	OtherExpense       float64 `json:"other_expense"`
	NetProfit          float64 `json:"net_profit"`
	ProfitMargin       float64 `json:"profit_margin"`
}

type ChartPointResponse struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Profit  float64 `json:"profit"`
}

type ExpenseBreakdownResponse struct {
	Label      string  `json:"label"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}