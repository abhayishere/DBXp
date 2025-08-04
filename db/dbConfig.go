package db

import (
	"fmt"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func (dc *DatabaseConfig) ConnectionString() string {
	_ = godotenv.Load() // Load .env file

	switch dc.Type {
	case "PostgreSQL":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dc.User, dc.Password, dc.Host, dc.Port, dc.Database)
	case "MySQL":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dc.User,
			dc.Password,
			dc.Host,
			dc.Port,
			dc.Database)
	default:
		return ""
	}
}

type QueryResult struct {
	Columns []string
	Rows    [][]string
}
