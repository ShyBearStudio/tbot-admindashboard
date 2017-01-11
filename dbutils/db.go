package dbutils

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDb(driver, dataSourceName string) (db *sql.DB, err error) {
	db, err = sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return
}
