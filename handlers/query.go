package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/abhayishere/DBXp/db"
	"github.com/rivo/tview"
)

type QueryHandler struct {
	db          db.Database
	resultBox   *tview.TextView
	refresh     func()
	history     *History
	export      *Export
	livePreview bool
	DebounceTimer *time.Timer
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
		livePreview: false, // Default to false, can be toggled later
		DebounceTimer: nil, // Initialize debounce timer
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
	if qh.IsLivePreviewEnabled() {
		if err == nil {
			qh.resultBox.SetText("[LIVE PREVIEW]\n" + output)
		}
	} else {
		qh.resultBox.SetText(output)
	}
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

func (qh *QueryHandler) IsLivePreviewEnabled() bool {
	return qh.livePreview
}
func (qh *QueryHandler) ShowLivePreview(query string) {
	if qh.IsSafeSelect(query) {
		start := time.Now()
		err := qh.ExecuteQuery(query)
		elapsed := time.Since(start)
		if err == nil {
			qh.resultBox.SetText(qh.resultBox.GetText(true) + "\nQuery executed in " + elapsed.String())
		}
	} else {
		qh.resultBox.SetText("Live preview only supports SELECT queries.")
	}
}

func (qh *QueryHandler) IsSafeSelect(query string) bool {
	sqlUpper := strings.ToUpper(strings.TrimSpace(query))
	if !strings.HasPrefix(sqlUpper, "SELECT") {
		return false
	}
	unsafeKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "TRUNCATE"}
	for _, keyword := range unsafeKeywords {
		if strings.Contains(sqlUpper, keyword) {
			return false
		}
	}
	return true
}
