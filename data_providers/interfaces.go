package data_providers

import (
    "database/sql"
)

type dataSource interface {
    Close()
    Exec(query string)(error)
    QueryRow(query string) *sql.Row
}
