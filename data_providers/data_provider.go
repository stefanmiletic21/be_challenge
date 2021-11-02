package data_providers

import (
    "database/sql"
    "fmt"
    "time"
)

type DataProvider struct {
    dataSource dataSource
}

func(dp *DataProvider) InitDataProvider(source dataSource) *DataProvider {
    dp.dataSource = source
    return dp
}

func NewDataProvider(source dataSource) *DataProvider {
    return (&DataProvider{}).InitDataProvider(source)
}

func(dp *DataProvider) MakeAggregatedTable() (err error) {
    err = dp.dataSource.Exec(`CREATE TABLE IF NOT EXISTS 
        aggregated_expenses (hour int PRIMARY KEY NOT NULL, expenses DOUBLE PRECISION)`)
    return err
}

func(dp *DataProvider) UpdateHour(hour int64) (err error) {
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

func(dp *DataProvider) GetFurthestUpdatedHour() (hour int64, err error) {
    err = dp.dataSource.QueryRow("SELECT MIN(hour) from aggregated_expenses").Scan(&hour)
    if err == sql.ErrNoRows {
        hour = time.Now().Truncate(time.Hour).Add(time.Hour).Unix()
        err = nil
        return
    }
    return
}