package db

import "database/sql"

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
	_, err := m.conn.Query(sql)
	if err != nil {
		return QueryResult{}, err
	}
	return QueryResult{}, nil
}

func (m *MySQLDB) ExecuteNonSelectQuery(sql string) (int64, error) {
	_, err := m.conn.Exec(sql)
	if err != nil {
		return 0, err
	}
	return 0, nil // Non-select queries don't return rows
}
