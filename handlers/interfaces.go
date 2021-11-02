package handlers

type APIServingService interface {
	GetLastHourExpenses() (expenses ExpensesByHourDTO, err error)
	GetExpenses(startHour, endHour int64) (expenses []ExpensesByHourDTO, err error)
}

type ExpensesByHourDTO struct {
	Timestamp int64   `json:"t"`
	Expenses  float64 `json:"v"`
}
