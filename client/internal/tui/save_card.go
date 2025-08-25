package tui

import (
	"strconv"

	"github.com/VicShved/pass-manager/client/internal/client"
	"github.com/rivo/tview"
)

func saveCard(app *tuiApplication, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.Box.SetBorder(true).SetTitle("Веедите логин/пароль")
	form.AddInputField("Card Number", "0000-0000-0000-0000", 19, nil, nil)
	form.AddInputField("Card Valid", "12/21", 5, nil, nil)
	form.AddInputField("Card Code", "000", 3, nil, nil)
	form.AddInputField("Description", "Desc about card", 50, nil, nil)
	form.AddButton("Cancel", func() {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	form.AddButton("Save", func() {
		cardNumber := form.GetFormItemByLabel("Card Number").(*tview.InputField).GetText()
		cardValid := form.GetFormItemByLabel("Card Valid").(*tview.InputField).GetText()
		cardCode := form.GetFormItemByLabel("Card Code").(*tview.InputField).GetText()
		cardDesc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
		_, rowID, err := client.DoSaveCard(app.tokenStr, cardNumber, cardValid, cardCode, cardDesc)
		modal := getModal(app, pages, "Успешная запись! Идентификатор = "+strconv.FormatUint(uint64(rowID), 10), "Ошибка записи! ", err)
		pages.AddAndSwitchToPage("Modal", modal, false)
		app.SetFocus(pages.SendToFront("mainMenuPage"))
	},
	)
	return form
}
