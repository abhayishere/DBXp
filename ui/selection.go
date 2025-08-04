package ui

import (
	"fmt"

	"github.com/abhayishere/DBXp/connection"
	"github.com/abhayishere/DBXp/db"
	"github.com/rivo/tview"
)

func NewConnectionSelector(app *tview.Application, onConnect func(db.DatabaseConfig)) {
	logo := tview.NewTextView()
	logo.SetText(`
 ██████╗ ██████╗ ██╗  ██╗██████╗ 
 ██╔══██╗██╔══██╗╚██╗██╔╝██╔══██╗
 ██║  ██║██████╔╝ ╚███╔╝ ██████╔╝
 ██║  ██║██╔══██╗ ██╔██╗ ██╔═══╝ 
 ██████╔╝██████╔╝██╔╝ ██╗██║     
 ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝     
                                 
     Database Explorer Tool      
`).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(false)
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Select Connection")
	list.AddItem("[1] Manual connection setup ⚙️", "", 0, func() {
		manualForm := NewManualConnectionForm(app, list, onConnect)
		app.SetRoot(manualForm, true)
	})
	list.AddItem("[2] Saved Connections", "", 0, func() {
		app.SetRoot(tview.NewTextView().SetText("Saved connections not implemented yet"), true)
	})
	list.AddItem("[3] Auto-detect local containers 🐳", "", 0, func() {
		list, err := connection.DetectDatabases(onConnect)
		if err != nil {
			app.SetRoot(tview.NewTextView().SetText("Detection from docker not implemented yet"), true)
		}
		if len(list) == 0 {
			app.SetRoot(tview.NewTextView().SetText("No databases detected"), true)
			return
		}
		listOfDbs := tview.NewList()
		listOfDbs.SetBorder(true).SetTitle("Detected Databases")
		for _, config := range list {
			listOfDbs.AddItem(fmt.Sprintf("%s of %s:%s", config.Type, config.Database, config.Port), "", 0, func() {
				onConnect(config)
			})
		}
		app.SetRoot(listOfDbs, true)
	})
	list.AddItem("[4] Set connection using env", "", 0, func() {
		app.SetRoot(tview.NewTextView().SetText("Set connection using env not implemented yet"), true)
	})
	list.AddItem("[5] Exit", "", 0, func() {
		app.Stop()
	})
	bottomBar := tview.NewTextView()

	bottomBar.SetText("[yellow][ ↑↓ ][white] navigate  [yellow][ enter ][white] select  [yellow][ esc ][white] back  [yellow][ ctrl+c ][white] exit").SetDynamicColors(true)
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(logo, 10, 0, false). // Logo takes 10 lines, fixed height
		AddItem(list, 0, 1, true).   // List takes remaining space
		AddItem(bottomBar, 1, 0, false)

	app.SetRoot(layout, true)
}

func NewManualConnectionForm(app *tview.Application, list *tview.List, onConnect func(db.DatabaseConfig)) *tview.Form {
	dbtype, host, port, user, password, database := "", "", "", "", "", ""
	form := tview.NewForm().
		AddDropDown("Database Type", []string{"PostgreSQL", "MySQL", "SQLite"}, 0, func(option string, optionIndex int) {
			dbtype = option
		}).
		AddInputField("Host", "", 20, nil, func(text string) {
			host = text
		}).
		AddInputField("Port", "", 5, nil, func(text string) {
			port = text
		}).
		AddInputField("User", "", 20, nil, func(text string) {
			user = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Database", "", 20, '*', func(text string) {
			database = text
		}).
		AddButton("Connect", func() {
			dbConfig := db.DatabaseConfig{
				Type:     dbtype,
				Host:     host,
				Port:     port,
				User:     user,
				Password: password,
				Database: database,
			}

			onConnect(dbConfig)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(list, true)
		})
	form.SetBorder(true).SetTitle("Manual Connection").SetTitleAlign(tview.AlignCenter)
	return form
}
