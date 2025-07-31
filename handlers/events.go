package handlers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// EventHandler manages UI event handling
type EventHandler struct {
	queryHandler *QueryHandler
}

// NewEventHandler creates a new event handler
func NewEventHandler(queryHandler *QueryHandler) *EventHandler {
	return &EventHandler{
		queryHandler: queryHandler,
	}
}

// SetupQueryInputHandler sets up the input field event handler
func (eh *EventHandler) SetupQueryInputHandler(queryInput *tview.InputField) {
	queryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			sql := queryInput.GetText()
			if sql != "" {
				eh.queryHandler.ExecuteQuery(sql)
				queryInput.SetText("") // Clear input after execution
			}
		}
	})
	queryInput.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyUp {
			previousQuery := eh.queryHandler.history.GetPreviousQuery()
			queryInput.SetText(previousQuery) // Set the previous query in the input field
			return nil
		} else if key.Key() == tcell.KeyDown {
			nextQuery := eh.queryHandler.history.GetNextQuery()
			queryInput.SetText(nextQuery) // Set the next query in the input field
			return nil
		} else if key.Key() == tcell.KeyCtrlE {
			err := eh.queryHandler.export.ExportToCSV() // Call export function when 'e' is pressed
			if err != nil {
				queryInput.SetText("Export Error: " + err.Error())
			} else {
				queryInput.SetText("Exported to export.csv") // Notify user of successful export
			}
			return nil
		}
		return key // Pass the event through
	})
}
