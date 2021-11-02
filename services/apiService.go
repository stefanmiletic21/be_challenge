package services

import (
	"fmt"
	"github.com/GlassNode/be_challenge/handlers"
)

type apiService struct {
	dataProvider APIDataProvider
}

func (ms *apiService) initAPIService(provider APIDataProvider) *apiService {
	ms.dataProvider = provider
	return ms
}

func NewAPIService(provider APIDataProvider) *apiService {
	return (&apiService{}).initAPIService(provider)
}

func (as *apiService) GetLastHourExpenses() (expenses handlers.ExpensesByHourDTO, err error) {
	expensesData, err := as.dataProvider.GetLastHourExpenses()
	if err != nil {
		fmt.Println(err)
		return
	}
	expenses = handlers.ExpensesByHourDTO(expensesData)
	return
}

func (as *apiService) GetExpenses(startHour, endHour int64) (expenses []handlers.ExpensesByHourDTO, err error) {
	expensesData, err := as.dataProvider.GetExpenses(startHour, endHour)
	if err != nil {
		fmt.Println(err)
		return
	}
	fetchedExpensesMap := make(map[int64]float64)
	for _, expense := range expensesData {
		fetchedExpensesMap[expense.Timestamp] = expense.Expenses
	}
	for i := startHour; i <= endHour; i = i + 3600 {
	    // if there is not aggregated_expense entry for some hour, we will return expense -1
		expenseVal := float64(-1)
		if val, ok := fetchedExpensesMap[i]; ok {
			expenseVal = val
		}
		expenses = append(expenses, handlers.ExpensesByHourDTO{
			Timestamp: i,
			Expenses:  expenseVal,
		})
	}
	return
}
