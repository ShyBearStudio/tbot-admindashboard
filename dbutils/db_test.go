package dbutils

import (
	"testing"
)

const (
	testDbDriver         = "postgres"
	testDbDataSourceName = "postgres://leabqlqclrgvcs:546b03954b61761cc2e3b2c0dbbfae83936db2fca407af6495ef41e2ada899aa@ec2-54-163-246-165.compute-1.amazonaws.com:5432/d7g3dt1edh1ta2"
)

func TestNewDbWithCorrectInputs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}
	if _, err := NewDb(testDbDriver, testDbDataSourceName); err != nil {
		t.Errorf("Cannot create new test database.")
	}
}

func TestNewDbWithIncorrectDataSource(t *testing.T) {
	if _, err := NewDb(testDbDriver, "Incorrect Data Source"); err == nil {
		t.Errorf("Should not create database connection due to incorrect data source name")
	}
}

func TestNewDbWithIncorrectDriver(t *testing.T) {
	if _, err := NewDb("IncorrectDriver", testDbDataSourceName); err == nil {
		t.Errorf("Should not create database connection due to incorrect driver")
	}
}
