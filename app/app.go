package app

import (
	"github.com/abhayishere/DBXp/contants"
	"github.com/abhayishere/DBXp/db"
	"github.com/abhayishere/DBXp/handlers"
	"github.com/abhayishere/DBXp/ui"
	"github.com/abhayishere/DBXp/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	tviewApp *tview.Application
}

func New() (*App, error) {
	return &App{
		tviewApp: tview.NewApplication(),
	}, nil
}

func (a *App) showConnectionSelector(onConnected func(db.Database)) (db.DatabaseConfig, error) {
	ui.NewConnectionSelector(a.tviewApp, func(dbConfig db.DatabaseConfig) {
		dbInstance, err := db.NewDatabase(dbConfig)
		if err != nil {
			a.tviewApp.SetRoot(tview.NewTextView().SetText("Failed to connect: "+err.Error()), true)
			return
		}
		onConnected(dbInstance)
	})
	return db.DatabaseConfig{}, nil
}
func (a *App) Run() error {
	_, err := a.showConnectionSelector(func(dbInstance db.Database) {
		a.showMainUILayout(dbInstance)
	})
	if err != nil {
		return err
	}
	return a.tviewApp.Run()
}

func (a *App) showMainUILayout(dbInstance db.Database) {
	resultBox := tview.NewTextView()
	resultBox.SetBorder(true).SetTitle("Results")

	queryInput := tview.NewInputField().SetLabel("SQL > ")
	queryInput.SetBorder(true)
	dbInstance.Connect()
	schemaList, refreshSchema := ui.GetSchemaExplorer(dbInstance, queryInput)
	queryHandler := handlers.NewQueryHandler(dbInstance, resultBox, refreshSchema)

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
	maincontent := ui.CreateLayoutWithHotKeys(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(utils.GetLogo(), 10, 0, false).
		AddItem(layout, 0, 1, true), contants.MainHotkeys)
	a.tviewApp.SetRoot(maincontent, true).SetFocus(queryInput)
}
