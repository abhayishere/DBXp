package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgreSQLDB struct {
	conn     *pgx.Conn
	dbConfig DatabaseConfig
}

func (p *PostgreSQLDB) Connect() error {
	var err error
	conString := p.dbConfig.ConnectionString()
	p.conn, err = pgx.Connect(context.Background(), conString)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSQLDB) ExecuteQuery(sql string) (QueryResult, error) {
	rows, err := p.conn.Query(context.Background(), sql)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name)
	}

	var results [][]string
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return QueryResult{}, err
		}
		rowData := make([]string, len(values))
		for i, value := range values {
			if value == nil {
				rowData[i] = "NULL" // Handle NULL values
			} else {
				rowData[i] = fmt.Sprintf("%v", value) // Convert other types to string
			}
		}
		results = append(results, rowData)
	}

	return QueryResult{
		Columns: columns,
		Rows:    results,
	}, nil
}

func (p *PostgreSQLDB) ExecuteNonSelectQuery(sql string) (int64, error) {
	_, err := p.conn.Exec(context.Background(), sql)
	if err != nil {
		return 0, err
	}
	return 0, nil // Non-select queries don't return rows
}
