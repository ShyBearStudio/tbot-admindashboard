package data

import (
	"database/sql"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb(driver, connectionString string) (err error) {
	logger.Tracef("Initializing database. Driver: '%s' | ConnectionString: '%s'", driver, connectionString)
	Db, err = sql.Open(driver, connectionString)
	if err != nil {
		logger.Errorln("Cannot open database", err)
	}
	return
}
