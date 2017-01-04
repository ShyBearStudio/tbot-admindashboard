package data

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

func InitDb(driver, connectionString string) (db *sql.DB, err error) {
	logger.Tracef("Initializing database. Driver: '%s' | ConnectionString: '%s'", driver, connectionString)
	db, err = sql.Open(driver, connectionString)
	if err != nil {
		logger.Errorln("Cannot open database", err)
	}
	logger.Tracef("Database was initialized successfully. Driver: '%s' | ConnectionString: '%s'", driver, connectionString)
	return
}

// Hashes plaintext with SHA-1
func Encrypt(textToEncrypt string) (encryptedText string) {
	encryptedText = fmt.Sprintf("%x", sha1.Sum([]byte(textToEncrypt)))
	return
}
