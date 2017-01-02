package data

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

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

// Hashes plaintext with SHA-1
func Encrypt(textToEncrypt string) (encryptedText string) {
	encryptedText = fmt.Sprintf("%x", sha1.Sum([]byte(textToEncrypt)))
	return
}
