package handlers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type EventHandler struct {
	queryHandler *QueryHandler
}

func NewEventHandler(queryHandler *QueryHandler) *EventHandler {
	return &EventHandler{
		queryHandler: queryHandler,
	}
}

func (eh *EventHandler) SetupQueryInputHandler(queryInput *tview.InputField) {
	queryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			sql := queryInput.GetText()
			if sql != "" {
				eh.queryHandler.ExecuteQuery(sql)
				queryInput.SetText("")
			}
		}
	})
	queryInput.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyUp {
			previousQuery := eh.queryHandler.history.GetPreviousQuery()
			queryInput.SetText(previousQuery)
			return nil
		} else if key.Key() == tcell.KeyDown {
			nextQuery := eh.queryHandler.history.GetNextQuery()
			queryInput.SetText(nextQuery)
			return nil
		} else if key.Key() == tcell.KeyCtrlE {
			err := eh.queryHandler.export.ExportToCSV()
			if err != nil {
				queryInput.SetText("Export Error: " + err.Error())
			} else {
				queryInput.SetText("Exported to export.csv")
			}
			return nil
		}
		return key
	})
}
