package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

type dbConnector struct {
	conn *sql.DB
}

func (db *dbConnector) initDB() (database *dbConnector, err error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPass := os.Getenv("POSTGRES_PASS")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return
	}

	psqlInfo := fmt.Sprintf("host=%v port=%d user=%v "+
		"password=%v dbname=%v sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}
	err = connection.Ping()
	if err != nil {
		return
	}

	db.conn = connection
	return db, nil
}

func (db *dbConnector) Close() {
	db.conn.Close()
}

func (db *dbConnector) Exec(query string) (err error) {
	_, err = db.conn.Exec(query)
	return err
}

func (db *dbConnector) QueryRow(query string) (row *sql.Row) {
	row = db.conn.QueryRow(query)
	return
}

func (db *dbConnector) Query(query string) (rows *sql.Rows, err error) {
	rows, err = db.conn.Query(query)
	return
}

func NewDBConnector() (db *dbConnector, err error) {
	db, err = (&dbConnector{}).initDB()
	return
}
