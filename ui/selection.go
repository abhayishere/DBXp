package ui

import (
	"fmt"

	"github.com/abhayishere/DBXp/connection"
	"github.com/abhayishere/DBXp/contants"
	"github.com/abhayishere/DBXp/db"
	"github.com/rivo/tview"
)

func NewConnectionSelector(app *tview.Application, onConnect func(db.DatabaseConfig)) {

	recreateMainMenu := func() {
		NewConnectionSelector(app, onConnect) // This is the right approach!
	}

	list := tview.NewList()
	list.SetBorder(true).SetTitle("Select Connection")
	list.AddItem("[1] Manual connection setup ⚙️", "", 0, func() {
		manualForm := NewManualConnectionForm(app, list, onConnect, recreateMainMenu)
		SetLayout(app, manualForm, contants.FormHotkeys, true, recreateMainMenu, recreateMainMenu)
	})
	list.AddItem("[2] Saved Connections", "", 0, func() {
		SetLayout(app, tview.NewTextView().SetText("Saved connections not implemented yet"), contants.FormHotkeys, false, recreateMainMenu, recreateMainMenu)
	})
	list.AddItem("[3] Auto-detect local containers 🐳", "", 0, func() {
		ListOfDatabases(app, onConnect, recreateMainMenu, recreateMainMenu)
	})
	list.AddItem("[4] Set connection using env", "", 0, func() {
		SetLayout(app, tview.NewTextView().SetText("Set connection using env not implemented yet"), contants.FormHotkeys, false, recreateMainMenu, recreateMainMenu)
	})
	list.AddItem("[5] Exit", "", 0, func() {
		app.Stop()
	})
	mainMenu := func() {
		SetLayout(app, list, contants.SelectorHotkeys, true, recreateMainMenu, recreateMainMenu)
	}
	mainMenu()
}

func NewManualConnectionForm(app *tview.Application, list *tview.List, onConnect func(db.DatabaseConfig), backScreen func()) *tview.Form {
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
			backScreen()
		})
	form.SetBorder(true).SetTitle("Manual Connection").SetTitleAlign(tview.AlignCenter)
	return form
}

func ListOfDatabases(app *tview.Application, onConnect func(db.DatabaseConfig), retry func(), backMenu func()) {
	list, err := connection.DetectDatabases(onConnect)
	if err != nil {
		SetLayout(app, tview.NewTextView().SetText("Error: "+err.Error()).SetTextAlign(tview.AlignCenter), contants.ErrorHotkeys, false, retry, backMenu)
		return
	}
	if len(list) == 0 {
		SetLayout(app, tview.NewTextView().SetText("Error: no databases detected").SetTextAlign(tview.AlignCenter), contants.ErrorHotkeys, false, retry, backMenu)
		return
	}
	listOfDbs := tview.NewList()
	listOfDbs.SetBorder(true).SetTitle("Detected Databases")
	for _, config := range list {
		listOfDbs.AddItem(fmt.Sprintf("%s of %s:%s", config.Type, config.Database, config.Port), "", 0, func() {
			onConnect(config)
		})
	}
	SetLayout(app, listOfDbs, contants.SelectorHotkeys, true, retry, backMenu)
}
