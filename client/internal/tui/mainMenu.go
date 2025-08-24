package tui

import "github.com/rivo/tview"

func mainMenu(app *tview.Application, pages *tview.Pages) tview.Primitive {
	list := tview.NewList()
	list.Box.SetBorder(true).SetTitle("Выберите действие")
	list.AddItem("Register", "", 'r', func() {
		pages.ShowPage("regLogPage")
		app.SetFocus(pages.SendToFront("regLogPage"))

	})
	list.AddItem("Exit", "", 'e', func() { app.Stop() })
	return list
}
