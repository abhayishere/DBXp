package utils

import "github.com/rivo/tview"

func GetLogo() *tview.TextView {
	logo := tview.NewTextView()
	logo.SetText(`
[cyan] ██████╗ ██████╗ ██╗  ██╗██████╗ [white]
[cyan] ██╔══██╗██╔══██╗╚██╗██╔╝██╔══██╗[white]
[cyan] ██║  ██║██████╔╝ ╚███╔╝ ██████╔╝[white]
[cyan] ██║  ██║██╔══██╗ ██╔██╗ ██╔═══╝ [white]
[cyan] ██████╔╝██████╔╝██╔╝ ██╗██║     [white]
[cyan] ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝     [white]
                                 
[yellow]     Database Explorer Tool      [white]
`).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(false)

	return logo
}
