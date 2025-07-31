package ui

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rivo/tview"
)

func GetSchemaExplorer(conn *pgx.Conn, queryInput *tview.InputField) (*tview.List, func()) {
	list := tview.NewList().ShowSecondaryText(false)

	refresh := func() {
		list.Clear()
		rows, _ := conn.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
		defer rows.Close()

		for rows.Next() {
			var tableName string
			_ = rows.Scan(&tableName)
			// Add item without callback in the AddItem call
			list.AddItem(tableName, "", 0, nil)
		}
	}
	refresh()

	// Set up the selection handler separately
	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// mainText contains the table name
		queryInput.SetText("SELECT * FROM " + mainText + ";")
	})

	list.SetBorder(true).SetTitle("Tables")
	return list, refresh
}
