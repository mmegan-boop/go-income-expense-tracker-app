package dto

type Response[T any] struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type CategorySummary struct {
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type SummaryResponse struct {
	TotalIncome       float64           `json:"total_income"`
	TotalExpense      float64           `json:"total_expense"`
	IncomeCategories  []CategorySummary `json:"income_categories"`
	ExpenseCategories []CategorySummary `json:"expense_categories"`
}
