package services

type WorkerDataProvider interface {
	MakeAggregatedTable() error
	UpdateHour(hour int64) (err error)
	GetFurthestUpdatedHour() (hour int64, err error)
}

type APIDataProvider interface {
	GetLastHourExpenses() (expense ExpensesByHour, err error)
	GetExpenses(startHour, endHour int64) (expenses []ExpensesByHour, err error)
}

type ExpensesByHour struct {
	Timestamp int64
	Expenses  float64
}
