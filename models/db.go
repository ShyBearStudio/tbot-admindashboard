package models

import (
	"database/sql"

	"github.com/ShyBearStudio/tbot-admindashboard/dbutils"
	_ "github.com/lib/pq"
)

type Datastore interface {
	CreateSession(*User) (Session, error)
	CheckSession(*Session) (bool, error)
	User(*Session) (User, error)
	UserByEmail(string) (User, error)
	AddUser(string, string, string, UserRoleType) (User, error)
	Users() ([]User, error)
}

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
