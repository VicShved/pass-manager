package tui

import (
	"strconv"

	"github.com/rivo/tview"
)

func getModal(app *tuiApplication, pages *tview.Pages, text string, textErr string, err error) *tview.Modal {
	modal := tview.NewModal()
	modalText := text
	if err != nil {
		modalText = textErr + err.Error()
	}
	modal.SetText(modalText)
	modal.AddButtons([]string{"Ok"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetFocus(pages.SendToFront("mainMenuPage"))
		pages.RemovePage("Modal")
	})
	return modal
}

func onlyUIntValue(textToCheck string, lastChar rune) bool {
	_, err := strconv.ParseUint(textToCheck, 10, 32)
	if err != nil {
		return false
	}
	return true
}
