package data

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb(driver, connectionString string) (err error) {
	Db, err = sql.Open(driver, connectionString)
	if err != nil {
		log.Fatal("Cannot open database", err)
	}
	return
}
