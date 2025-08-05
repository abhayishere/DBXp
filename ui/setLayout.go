package ui

import (
	"github.com/abhayishere/DBXp/contants"
	"github.com/abhayishere/DBXp/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetLayout(app *tview.Application, content tview.Primitive, hotkeys string, success bool, retry func(), mainMenu func()) {
	var maincontent *tview.Flex
	if success {
		app.SetInputCapture(nil)
		maincontent = CreateLayoutWithHotKeys(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(utils.GetLogo(), 10, 0, false).
			AddItem(content, 0, 1, true), hotkeys)

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEsc:
				retry()
				return nil
			case tcell.KeyCtrlC:
				app.Stop()
				return nil
			}
			return event
		})
		app.SetRoot(maincontent, true)
	} else {
		maincontent = CreateLayoutWithHotKeys(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(utils.GetLogo(), 10, 0, false).
			AddItem(content, 0, 1, true), contants.ErrorHotkeys)

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEsc:
				mainMenu()
				return nil
			case tcell.KeyCtrlC:
				app.Stop()
				return nil
			case tcell.KeyEnter:
				retry()
				return nil
			}
			return event
		})
		app.SetRoot(maincontent, true)
	}
}
