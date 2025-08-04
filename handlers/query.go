package handlers

import (
	"fmt"
	"strings"

	"github.com/abhayishere/DBXp/db"
	"github.com/rivo/tview"
)

type QueryHandler struct {
	db        db.Database
	resultBox *tview.TextView
	refresh   func()
	history   *History
	export    *Export
}

func NewQueryHandler(db db.Database, resultBox *tview.TextView, refreshFunc func()) *QueryHandler {
	return &QueryHandler{
		db:        db,
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

func (qh *QueryHandler) ExecuteQuery(sql string) error {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if len(sqlUpper) != 0 {
		qh.history.history = append(qh.history.history, sqlUpper)
		qh.history.historyindex++
	}
	if strings.HasPrefix(sqlUpper, "SELECT") {
		return qh.executeSelectQuery(sql)
	} else {
		return qh.executeNonSelectQuery(sql, sqlUpper)
	}
}

func (qh *QueryHandler) executeSelectQuery(sql string) error {
	queryResult, err := qh.db.ExecuteQuery(sql)
	if err != nil {
		qh.resultBox.SetText("Error: " + err.Error())
		return err
	}
	output := qh.formatQueryResults(queryResult)
	qh.resultBox.SetText(output)
	return nil
}

func (qh *QueryHandler) executeNonSelectQuery(sql, sqlUpper string) error {
	_, err := qh.db.ExecuteNonSelectQuery(sql)
	if err != nil {
		qh.resultBox.SetText("Error: " + err.Error())
		return err
	}

	qh.resultBox.SetText("Query executed successfully")

	if qh.shouldRefreshSchema(sqlUpper) {
		qh.refresh()
	}
	return nil
}

func (qh *QueryHandler) formatQueryResults(data db.QueryResult) string {
	var output string

	fields := data.Columns
	qh.export.AddColumns(fields)
	for _, f := range fields {
		output += f + "\t"
	}
	output += "\n"

	for _, row := range data.Rows {
		for _, value := range row {
			if value != "NULL" && value != "" {
				output += fmt.Sprintf("%v", value) + "\t"
			} else {
				output += "NULL\t"
			}
		}
		output += "\n"
		qh.export.AddRow(row)
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
