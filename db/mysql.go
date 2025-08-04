package db

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	conn     *sql.DB
	dbConfig DatabaseConfig
}

func (m *MySQLDB) Connect() error {
	var err error
	m.conn, err = sql.Open("mysql", m.dbConfig.ConnectionString())
	if err != nil {
		return err
	}
	return m.conn.Ping()
}

func (m *MySQLDB) ExecuteQuery(sql string) (QueryResult, error) {
	rows, err := m.conn.Query(sql)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()

	fieldDescriptions, err := rows.ColumnTypes()
	if err != nil {
		return QueryResult{}, err
	}
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name()
	}
	var results [][]string
	for rows.Next() {
		values := make([]string, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range columns {
			scanArgs[i] = &values[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return QueryResult{}, err
		}
		results = append(results, values)
	}
	return QueryResult{Columns: columns, Rows: results}, nil
}

func (m *MySQLDB) ExecuteNonSelectQuery(sql string) (int64, error) {
	_, err := m.conn.Exec(sql)
	if err != nil {
		return 0, errors.New("Error executing non-select query: " + err.Error())
	}
	return 0, nil // Non-select queries don't return rows
}

func (m *MySQLDB) ListTables() ([]string, error) {
	rows, err := m.conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}
