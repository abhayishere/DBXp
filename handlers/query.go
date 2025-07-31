package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rivo/tview"
)

// QueryHandler handles SQL query execution and result formatting
type QueryHandler struct {
	conn      *pgx.Conn
	resultBox *tview.TextView
	refresh   func()
	history   *History // History of executed queries
	export    *Export  // Export handler for exporting results
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(conn *pgx.Conn, resultBox *tview.TextView, refreshFunc func()) *QueryHandler {
	return &QueryHandler{
		conn:      conn,
		resultBox: resultBox,
		refresh:   refreshFunc,
		history: &History{
			history:      []string{},
			historyindex: -1, // Start with no history index
		},
		export: &Export{
			columns: []string{},
			rows:    [][]string{},
		},
	}
}

// ExecuteQuery handles both SELECT and non-SELECT queries
func (qh *QueryHandler) ExecuteQuery(sql string) {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if len(sqlUpper) != 0 {
		qh.history.history = append(qh.history.history, sqlUpper) // Store the executed query in history
		qh.history.historyindex++                                 // Update history index to the latest query
	}
	if strings.HasPrefix(sqlUpper, "SELECT") {
		qh.executeSelectQuery(sql)
	} else {
		qh.executeNonSelectQuery(sql, sqlUpper)
	}
}

// executeSelectQuery handles SELECT queries and formats results
func (qh *QueryHandler) executeSelectQuery(sql string) {
	rows, err := qh.conn.Query(context.Background(), sql)
	if err != nil {
		qh.resultBox.SetText("Error: " + err.Error())
		return
	}
	defer rows.Close()

	output := qh.formatQueryResults(rows)
	qh.resultBox.SetText(output)
}

// executeNonSelectQuery handles INSERT, UPDATE, DELETE, CREATE, DROP, etc.
func (qh *QueryHandler) executeNonSelectQuery(sql, sqlUpper string) {
	_, err := qh.conn.Exec(context.Background(), sql)
	if err != nil {
		qh.resultBox.SetText("Error: " + err.Error())
		return
	}

	qh.resultBox.SetText("Query executed successfully")

	// Refresh schema if it's a DDL operation that might affect table structure
	if qh.shouldRefreshSchema(sqlUpper) {
		qh.refresh()
	}
}

// formatQueryResults formats the query results into a readable string
func (qh *QueryHandler) formatQueryResults(rows pgx.Rows) string {
	var output string

	// Add column headers
	fields := rows.FieldDescriptions()
	qh.export.AddColumns(fields)
	for _, f := range fields {
		output += f.Name + "\t"
	}
	output += "\n"

	// Add data rows
	for rows.Next() {
		values, _ := rows.Values()
		for _, v := range values {
			if v != nil {
				output += fmt.Sprintf("%v", v) + "\t"
			} else {
				output += "NULL\t"
			}
		}
		output += "\n"
		qh.export.AddRow(values) // Store row data for export
	}

	return output
}

// shouldRefreshSchema determines if schema should be refreshed based on SQL command
func (qh *QueryHandler) shouldRefreshSchema(sqlUpper string) bool {
	ddlCommands := []string{"CREATE", "DROP", "ALTER", "TRUNCATE"}
	for _, cmd := range ddlCommands {
		if strings.HasPrefix(sqlUpper, cmd) {
			return true
		}
	}
	return false
}
