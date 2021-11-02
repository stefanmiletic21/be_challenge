package data_providers

import (
	"database/sql"
	"fmt"
	"github.com/GlassNode/be_challenge/services"
	"time"
)

type DataProvider struct {
	dataSource dataSource
}

func (dp *DataProvider) InitDataProvider(source dataSource) *DataProvider {
	dp.dataSource = source
	return dp
}

func NewDataProvider(source dataSource) *DataProvider {
	return (&DataProvider{}).InitDataProvider(source)
}

func (dp *DataProvider) MakeAggregatedTable() (err error) {
	err = dp.dataSource.Exec(`CREATE TABLE IF NOT EXISTS 
        aggregated_expenses (hour int PRIMARY KEY NOT NULL, expenses DOUBLE PRECISION)`)
	return err
}

func (dp *DataProvider) UpdateHour(hour int64) (err error) {
	previousHour := hour - 3600
	err = dp.dataSource.Exec(fmt.Sprintf(`INSERT INTO aggregated_expenses(hour,expenses) VALUES(%d,
                                        (select sum(gas_price*gas_used/1000000000/1000000000)
                                            from public.transactions
                                            where block_time BETWEEN to_timestamp(%d) AND to_timestamp(%d)
                                                and "from" not in (select address
                                                    from public.contracts)
                                                and "to" not in (select address
                                                    from public.contracts)
                                                    and "from" != '0x000000000000000000000000000000000000000'
                                                    and "to" != '0x000000000000000000000000000000000000000'))
                                            ON CONFLICT(hour) DO UPDATE SET expenses = EXCLUDED.expenses`,
		hour, previousHour, hour))
	return err
}

func (dp *DataProvider) GetFurthestUpdatedHour() (hour int64, err error) {
	err = dp.dataSource.QueryRow("SELECT MIN(hour) from aggregated_expenses").Scan(&hour)
	if err == sql.ErrNoRows {
		hour = time.Now().Truncate(time.Hour).Add(time.Hour).Unix()
		err = nil
		return
	}
	return
}

func (dp *DataProvider) GetLastHourExpenses() (expense services.ExpensesByHour, err error) {
	hour := time.Now().Truncate(time.Hour).Add(time.Hour).Unix()
	var expenseVal sql.NullFloat64
	err = dp.dataSource.QueryRow(fmt.Sprintf("SELECT expenses FROM aggregated_expenses WHERE hour = %d", hour)).Scan(&expenseVal)
	if err != nil {
		return
	}
	expense.Timestamp = hour
	if expenseVal.Valid {
		expense.Expenses = expenseVal.Float64
	} else {
		expense.Expenses = -1
	}
	return
}

func (dp *DataProvider) GetExpenses(startHour, endHour int64) (expenses []services.ExpensesByHour, err error) {
	rows, err := dp.dataSource.Query(fmt.Sprintf("SELECT * FROM aggregated_expenses WHERE hour >= %d AND hour <= %d", startHour, endHour))
	if err != nil {
		return
	}

	for rows.Next() {
		var dbHour sql.NullInt64
		var dbExpenses sql.NullFloat64
		if err := rows.Scan(&dbHour, &dbExpenses); err != nil {
			return expenses, err
		}
		exp := services.ExpensesByHour{
			Timestamp: dbHour.Int64,
		}
		if dbExpenses.Valid {
			exp.Expenses = dbExpenses.Float64
		} else {
			exp.Expenses = -1
		}
		expenses = append(expenses, exp)
	}
	return
}
