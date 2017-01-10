package models

import (
	"database/sql"

	"github.com/ShyBearStudio/tbot-admindashboard/dbutils"
	_ "github.com/lib/pq"
)

type Db struct {
	*sql.DB
}

func NewDb(driver, dataSourceName string) (*Db, error) {
	db, err := dbutils.NewDb(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Db{db}, nil
}
