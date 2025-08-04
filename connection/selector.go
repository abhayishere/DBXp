package connection

import (
	"github.com/rivo/tview"
)

type ConnectionSelector struct {
	app     *tview.Application
	options []ConnectionOption
}

type ConnectionOption struct {
	Name        string
	Type        string
	Description string
}

