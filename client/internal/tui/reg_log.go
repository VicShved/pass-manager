package tui

import (
	"github.com/VicShved/pass-manager/client/internal/client"
	"github.com/rivo/tview"
)

func regLog(app *tview.Application, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Веедите логин/пароль")
	form.AddInputField("Login", "", 20, nil, nil)
	form.AddInputField("Password", "enter password", 20, nil, nil)
	form.AddButton("Register", func() {
		login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
		pswrd := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		_, _, err := client.DoRegister(login, pswrd)
		if err != nil {
			panic(err)
		}
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
