package tui

import "github.com/rivo/tview"

// StartTui run text user innterface
func StartTui() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	mainMenuPage := mainMenu(app, pages)
	regLogPage := regLog(app, pages)
	pages.AddPage("mainMenuPage", mainMenuPage, true, true)
	pages.AddPage("regLogPage", regLogPage, true, false)
	if err := app.SetRoot(pages, true).SetFocus(mainMenuPage).Run(); err != nil {
		panic(err)
	}
}
