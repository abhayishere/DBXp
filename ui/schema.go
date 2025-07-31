package ui

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rivo/tview"
)

func GetSchemaExplorer(conn *pgx.Conn, queryInput *tview.InputField) (*tview.List, func()) {
	list := tview.NewList().ShowSecondaryText(false)

	refresh := func() {
		tableList := []string{}
		list.Clear()
		rows, _ := conn.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
		defer rows.Close()
		for rows.Next() {
			var tableName string
			_ = rows.Scan(&tableName)
			tableList = append(tableList, tableName)
		}
		rows.Close()
		for _, table_name := range tableList {
			rows, _ := conn.Query(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM %s", table_name))
			defer rows.Close()
			var count int64
			if rows.Next() {
				_ = rows.Scan(&count)
			}
			list.AddItem(fmt.Sprintf("%s(%d)", table_name, count), "", 0, nil)
		}
	}
	refresh()

	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		tableName := mainText
		for i, char := range mainText {
			if char == '(' {
				tableName = mainText[:i]
				break
			}
		}
		queryInput.SetText("SELECT * FROM " + tableName + ";")
	})

	list.SetBorder(true).SetTitle("Tables")
	return list, refresh
}
