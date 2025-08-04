package ui

import (
	"github.com/rivo/tview"
)

func BuildUILayout(schemaList *tview.List, queryInput *tview.InputField, resultBox *tview.TextView) *tview.Flex {

	mainContent := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(resultBox, 0, 4, false).
		AddItem(queryInput, 3, 1, false)

	root := tview.NewFlex().AddItem(schemaList, 30, 1, false).AddItem(mainContent, 0, 3, true)

	return root
}
