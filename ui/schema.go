package ui

import (
	"fmt"

	"github.com/abhayishere/DBXp/db"
	"github.com/rivo/tview"
)

func GetSchemaExplorer(db db.Database, queryInput *tview.InputField) (*tview.List, func()) {
	list := tview.NewList().ShowSecondaryText(false)
	refresh := func() {
		list.Clear()
		queryResult, _ := db.ExecuteQuery("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
		tableList := queryResult.Rows
		for _, table_name := range tableList {
			queryResult, _ := db.ExecuteQuery(fmt.Sprintf("SELECT COUNT(*) FROM %s", table_name[0]))
			list.AddItem(fmt.Sprintf("%s(%s)", table_name[0], queryResult.Rows[0][0]), "", 0, nil)
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
