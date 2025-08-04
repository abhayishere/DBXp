package ui

import (
	"github.com/abhayishere/DBXp/db"
	"github.com/rivo/tview"
)

func NewConnectionSelector(app *tview.Application, onConnect func(db.DatabaseConfig)) {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Select Connection")
	list.AddItem("Manual Connection", "Configure a connection manually", 1, func() {
		manualForm := NewManualConnectionForm(app, list, onConnect)
		app.SetRoot(manualForm, true)
	})
	list.AddItem("Saved Connections", "Select from saved connections", 1, func() {
		app.SetRoot(tview.NewTextView().SetText("Saved connections not implemented yet"), true)
	})
	list.AddItem("Set connection using env", "Update the env file in the dir", 1, func() {
		app.SetRoot(tview.NewTextView().SetText("Set connection using env not implemented yet"), true)
	})
	list.AddItem("Exit", "Exit the connection selector", 1, func() {
		app.Stop()
	})
	app.SetRoot(list, true)
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
