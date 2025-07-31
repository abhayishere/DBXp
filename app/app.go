package app

import (
	"context"

	"github.com/abhayishere/DBXp/db"
	"github.com/abhayishere/DBXp/handlers"
	"github.com/abhayishere/DBXp/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rivo/tview"
)

type App struct {
	tviewApp *tview.Application
	conn     *pgx.Conn
}

func New() (*App, error) {
	conn, err := db.Connect()
	if err != nil {
		return nil, err
	}

	return &App{
		tviewApp: tview.NewApplication(),
		conn:     conn,
	}, nil
}

func (a *App) Run() error {
	defer a.conn.Close(context.Background())

	resultBox := tview.NewTextView()
	resultBox.SetBorder(true).SetTitle("Results")

	queryInput := tview.NewInputField().SetLabel("SQL > ")

	schemaList, refreshSchema := ui.GetSchemaExplorer(a.conn, queryInput)

	queryHandler := handlers.NewQueryHandler(a.conn, resultBox, refreshSchema)

	eventHandler := handlers.NewEventHandler(queryHandler)
	eventHandler.SetupQueryInputHandler(queryInput)

	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:

			if a.tviewApp.GetFocus() == schemaList {
				a.tviewApp.SetFocus(queryInput)
			} else {
				a.tviewApp.SetFocus(schemaList)
			}
			return nil
		}
		return event
	})

	layout := ui.BuildUILayout(schemaList, queryInput, resultBox)

	a.tviewApp.SetFocus(queryInput)

	return a.tviewApp.SetRoot(layout, true).Run()
}
