package ui

import (
	"github.com/rivo/tview"
)

func CreateLayoutWithHotKeys(content tview.Primitive, hotkeysText string) *tview.Flex {
	bottomBar := tview.NewTextView()
	bottomBar.SetText(hotkeysText).SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(content, 0, 1, true).
		AddItem(bottomBar, 3, 0, false)
	return layout
}
