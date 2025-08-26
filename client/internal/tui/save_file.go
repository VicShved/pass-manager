package tui

import (
	"strconv"

	"github.com/rivo/tview"
)

func saveFile(app *tuiApplication, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Веедите логин/пароль")
	form.AddInputField("Filename", "", 20, nil, nil)
	form.AddInputField("Description", "Desc about file", 50, nil, nil)
	form.AddButton("Cancel", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Save", func() {
		fileName := form.GetFormItemByLabel("Filename").(*tview.InputField).GetText()
		cardDesc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
		_, rowID, err := app.client.DoSaveFile(app.tokenStr, fileName, cardDesc)
		modal := getModal(app, pages, "Успешная запись! Идентификатор = "+strconv.FormatUint(uint64(rowID), 10), "Ошибка записи! ", err)
		pages.AddAndSwitchToPage("Modal", modal, false)
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	return form
}
