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

// App represents the main application
type App struct {
	tviewApp *tview.Application
	conn     *pgx.Conn
}

// New creates a new application instance
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

// Run starts the application
func (a *App) Run() error {
	defer a.conn.Close(context.Background())

	// Create UI components
	resultBox := tview.NewTextView()
	resultBox.SetBorder(true).SetTitle("Results")

	queryInput := tview.NewInputField().SetLabel("SQL > ")

	schemaList, refreshSchema := ui.GetSchemaExplorer(a.conn, queryInput)

	// Set up handlers
	queryHandler := handlers.NewQueryHandler(a.conn, resultBox, refreshSchema)
	
	eventHandler := handlers.NewEventHandler(queryHandler)
	eventHandler.SetupQueryInputHandler(queryInput)

	// Set up global input capture for Tab navigation
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Toggle focus between schema list and query input
			if a.tviewApp.GetFocus() == schemaList {
				a.tviewApp.SetFocus(queryInput)
			} else {
				a.tviewApp.SetFocus(schemaList)
			}
			return nil // Consume the event
		}
		return event // Pass through other events
	})

	// Build layout
	layout := ui.BuildUILayout(schemaList, queryInput, resultBox)

	// Set initial focus to query input
	a.tviewApp.SetFocus(queryInput)

	// Run the application
	return a.tviewApp.SetRoot(layout, true).Run()
}
