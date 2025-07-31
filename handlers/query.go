package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rivo/tview"
)

type QueryHandler struct {
	conn      *pgx.Conn
	resultBox *tview.TextView
	refresh   func()
	history   *History
	export    *Export
}

func NewQueryHandler(conn *pgx.Conn, resultBox *tview.TextView, refreshFunc func()) *QueryHandler {
	return &QueryHandler{
		conn:      conn,
		resultBox: resultBox,
		refresh:   refreshFunc,
		history: &History{
			history:      []string{},
			historyindex: -1,
		},
		export: &Export{
			columns: []string{},
			rows:    [][]string{},
		},
	}
}

func (qh *QueryHandler) ExecuteQuery(sql string) {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if len(sqlUpper) != 0 {
		qh.history.history = append(qh.history.history, sqlUpper)
		qh.history.historyindex++
	}
	if strings.HasPrefix(sqlUpper, "SELECT") {
		qh.executeSelectQuery(sql)
	} else {
		qh.executeNonSelectQuery(sql, sqlUpper)
	}
}

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

func (qh *QueryHandler) executeNonSelectQuery(sql, sqlUpper string) {
	_, err := qh.conn.Exec(context.Background(), sql)
	if err != nil {
		qh.resultBox.SetText("Error: " + err.Error())
		return
	}

	qh.resultBox.SetText("Query executed successfully")

	if qh.shouldRefreshSchema(sqlUpper) {
		qh.refresh()
	}
}

func (qh *QueryHandler) formatQueryResults(rows pgx.Rows) string {
	var output string

	fields := rows.FieldDescriptions()
	qh.export.AddColumns(fields)
	for _, f := range fields {
		output += f.Name + "\t"
	}
	output += "\n"

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
		qh.export.AddRow(values)
	}

	return output
}

func (qh *QueryHandler) shouldRefreshSchema(sqlUpper string) bool {
	ddlCommands := []string{"CREATE", "DROP", "ALTER", "TRUNCATE"}
	for _, cmd := range ddlCommands {
		if strings.HasPrefix(sqlUpper, cmd) {
			return true
		}
	}
	return false
}
