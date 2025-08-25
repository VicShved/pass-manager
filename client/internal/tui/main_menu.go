package tui

import "github.com/rivo/tview"

func mainMenu(app *tuiApplication, pages *tview.Pages) tview.Primitive {
	list := tview.NewList()
	list.Box.SetBorder(true).SetTitle("Выберите действие")
	list.AddItem("Register or Login", "", 'r', func() {
		pages.ShowPage("regLogPage")
		app.SetFocus(pages.SendToFront("regLogPage"))
	})
	list.AddItem("---------------------------------------------------------------------", "", rune(0), nil)
	list.AddItem("Save Login & Password", "", 'l', func() {
		pages.ShowPage("saveLogPass")
		app.SetFocus(pages.SendToFront("saveLogPass"))
	})
	list.AddItem("Get Login & Password", "", 'm', func() {
		pages.ShowPage("getLogPass")
		app.SetFocus(pages.SendToFront("getLogPass"))
	})
	list.AddItem("---------------------------------------------------------------------", "", rune(0), nil)
	list.AddItem("Save Card", "", 'c', func() {
		pages.ShowPage("saveCard")
		app.SetFocus(pages.SendToFront("saveCard"))
	})
	list.AddItem("Get Card", "", 'd', func() {
		pages.ShowPage("getCard")
		app.SetFocus(pages.SendToFront("getCard"))
	})
	list.AddItem("---------------------------------------------------------------------", "", rune(0), nil)
	list.AddItem("Save File", "", 'f', func() {
		pages.ShowPage("saveFile")
		app.SetFocus(pages.SendToFront("saveFile"))
	})
	list.AddItem("Get File", "", 'f', func() {
		pages.ShowPage("getFile")
		app.SetFocus(pages.SendToFront("getFile"))
	})
	list.AddItem("---------------------------------------------------------------------", "", rune(0), nil)
	list.AddItem("Exit", "", 'e', func() { app.Stop() })
	return list
}
