package dbutils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	scriptsRepo string = `..\sql`
)

func ExecuteSqlScript(db *sql.DB, scriptName string) error {
	fullPath := filepath.Join(scriptsRepo, scriptName)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	commands := strings.Split(string(content), ";")
	for _, command := range commands {
		command = strings.Trim(command, " \t\n\r\n")
		if len(command) == 0 {
			continue
		}
		fmt.Printf("Executing: '%v' ...\n", command)
		_, err := db.Exec(command)
		if err != nil {
			return err
		}
	}
	return nil
}
