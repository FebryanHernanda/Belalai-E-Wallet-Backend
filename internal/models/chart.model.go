package models

type ChartData struct {
	Labels      []string `json:"labels"`
	IncomeData  []int    `json:"income_data"`
	ExpenseData []int    `json:"expense_data"`
}

type ChartDataResponse struct {
	Response
	Data ChartData `json:"data"`
}
