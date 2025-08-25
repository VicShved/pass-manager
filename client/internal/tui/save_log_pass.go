package tui

import (
	"strconv"

	"github.com/VicShved/pass-manager/client/internal/client"
	"github.com/rivo/tview"
)

func saveLogPass(app *tuiApplication, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Веедите логин/пароль")
	form.AddInputField("Login", "", 20, nil, nil)
	form.AddInputField("Password", "", 20, nil, nil)
	form.AddInputField("Description", "Desc about login/password", 50, nil, nil)
	form.AddButton("Cancel", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Save", func() {
		login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		desc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
		_, rowID, err := client.DoSaveLogPass(app.tokenStr, login, password, desc)
		modal := getModal(app, pages, "Успешная запись! Идентификатор = "+strconv.FormatUint(uint64(rowID), 10), "Ошибка записи! ", err)
		pages.AddAndSwitchToPage("Modal", modal, false)
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	return form
}
