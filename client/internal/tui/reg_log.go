package tui

import (
	"github.com/rivo/tview"
)

func regLog(app *tuiApplication, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Веедите логин/пароль")
	form.AddInputField("Login", "", 20, nil, nil)
	form.AddInputField("Password", "enter password", 20, nil, nil)
	form.AddButton("Login", func() {
		login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
		pswrd := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		_, tokenStr, err := app.client.DoLogin(login, pswrd)
		if err == nil {
			app.tokenStr = tokenStr
		}
		modal := getModal(app, pages, "Успешная авторизация!", "Ошибка авторизации! ", err)
		pages.AddAndSwitchToPage("Modal", modal, false)
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Cancel", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Register", func() {
		login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
		pswrd := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		_, tokenStr, err := app.client.DoRegister(login, pswrd)
		if err == nil {
			app.tokenStr = tokenStr
		}
		modal := getModal(app, pages, "Успешная регистрация!", "Ошибка регистрации! ", err)
		pages.AddAndSwitchToPage("Modal", modal, false)
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	return form
}
