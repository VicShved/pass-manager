package tui

import (
	"strconv"

	"github.com/rivo/tview"
)

func getCard(app *tuiApplication, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Загрузка данных карты")
	form.AddInputField("Row ID", "", 10, onlyUIntValue, nil)
	form.AddButton("Get", func() {
		rowIDStr := form.GetFormItemByLabel("Row ID").(*tview.InputField).GetText()
		rowID, err := strconv.ParseUint(rowIDStr, 10, 32)
		var modal *tview.Modal
		if err != nil {
			modal = getModal(app, pages, "", "Ошибка конвертации '"+rowIDStr+"' в целое число. Введите корректное число!", err)

		} else {
			_, dataStr, err := app.client.DoGetCard(app.tokenStr, uint32(rowID))
			modal = getModal(app, pages, "Полученные данные!\n"+dataStr, "Данные с таким ID отсутствуют! ", err)
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
