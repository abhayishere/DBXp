package db

import (
	"fmt"
)

type Database interface {
	// Connect establishes a connection to the database.
	Connect() error
	// ExecuteQuery executes a SQL query and returns the results.
	ExecuteQuery(sql string) (QueryResult, error)
	// GetTables retrieves a list of tables in the database.
	ExecuteNonSelectQuery(sql string) (int64, error)

	ListTables() ([]string, error)
}

func NewDatabase(dbConfig DatabaseConfig) (Database, error) {
	if dbConfig.Type == "PostgreSQL" {
		return &PostgreSQLDB{dbConfig: dbConfig}, nil
	}
	if dbConfig.Type == "MySQL" {
		return &MySQLDB{dbConfig: dbConfig}, nil
	}
	return nil, fmt.Errorf("unsupported database type: %s", dbConfig.Type)
}
