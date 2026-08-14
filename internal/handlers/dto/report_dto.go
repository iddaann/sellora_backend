package dto 

type SummaryResponse struct {
	Revenue				float64 `json:"revenue"`
	COGS				float64 `json:"cogs"`
	GrossProfit 		float64 `json:"gross_profit"`
	OperationalExpense	float64 `json:"operational_expense"`
	OtherExpense		float64 `json:"other_expense"`
	NetPofit			float64 `json:"net_profit"`
	ProfitMargin		float64 `json:"profit_margin"`
}

type ChartPointResponse struct {
	Date 	string 	`json:"date"`
	Revenue float64 `json:"revenue"`
	Profit 	float64	`json:"profit"`
}

type ExpenseBreakdownResponse struct {
	Label		string 	`json:"abel"`
	Amount 		float64 `json:"amount"`
	Percentage	float64 `json:"percenetage"`
}