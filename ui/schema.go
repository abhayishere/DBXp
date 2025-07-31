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

			list.AddItem(tableName, "", 0, nil)
		}
	}
	refresh()

	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {

		queryInput.SetText("SELECT * FROM " + mainText + ";")
	})

	list.SetBorder(true).SetTitle("Tables")
	return list, refresh
}
