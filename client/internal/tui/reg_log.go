package tui

import (
	"github.com/rivo/tview"
)

func regLog(app *tview.Application, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Веедите логин/пароль")
	form.AddInputField("Login", "enter login", 20, nil, nil)
	form.AddInputField("Password", "enter password", 20, nil, nil)
	form.AddButton("Register", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Cancel", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Login", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)

	// modal.AddButtons([]string{"Quit", "Cancel"})
	// modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	// 	if buttonLabel == "Quit" {
	// 		app.Stop()
	// 	}
	// })
	return form
}
