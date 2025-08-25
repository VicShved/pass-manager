package tui

import (
	"strconv"

	"github.com/VicShved/pass-manager/client/internal/client"
	"github.com/rivo/tview"
)

func getFile(app *tuiApplication, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Загрузка файла")
	form.AddInputField("Row ID", "", 10, onlyUIntValue, nil)
	form.AddButton("Load", func() {
		rowIDStr := form.GetFormItemByLabel("Row ID").(*tview.InputField).GetText()
		rowID, err := strconv.ParseUint(rowIDStr, 10, 32)
		var modal *tview.Modal
		if err != nil {
			modal = getModal(app, pages, "", "Ошибка конвертации '"+rowIDStr+"' в целое число. Введите корректное число!", err)

		} else {
			_, fileName, err := client.DoGetFile(app.tokenStr, uint32(rowID))
			modal = getModal(app, pages, "Файл  записан в\n"+fileName, "Данные с таким ID отсутствуют! ", err)
		}

		pages.AddAndSwitchToPage("Modal", modal, false)
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Cancel", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	return form
}
