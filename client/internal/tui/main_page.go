package tui

import (
	"github.com/VicShved/pass-manager/client/internal/client"
	"github.com/VicShved/pass-manager/client/internal/config"
	"github.com/rivo/tview"
)

type tuiApplication struct {
	tview.Application
	// serverUrl string
	tokenStr  string
	client    *client.GClient
	version   string
	buildTime string
}

// StartTui run text user innterface
func StartTui(version string, buildTime string) {
	conf := config.GetClientConfig()
	gClient := client.GetgClient(*conf)
	app := tuiApplication{
		Application: *tview.NewApplication(),
		client:      gClient,
		version:     version,
		buildTime:   buildTime,
	}
	pages := tview.NewPages()
	mainMenuPage := mainMenu(&app, pages)
	regLogPage := regLog(&app, pages)
	saveLogPass := saveLogPass(&app, pages)
	getLogPass := getLogpass(&app, pages)
	saveCard := saveCard(&app, pages)
	getCard := getCard(&app, pages)
	saveFile := saveFile(&app, pages)
	getFile := getFile(&app, pages)
	pages.AddPage("mainMenuPage", mainMenuPage, true, true)
	pages.AddPage("regLogPage", regLogPage, true, false)
	pages.AddPage("saveLogPass", saveLogPass, true, false)
	pages.AddPage("getLogPass", getLogPass, true, false)
	pages.AddPage("saveCard", saveCard, true, false)
	pages.AddPage("getCard", getCard, true, false)
	pages.AddPage("saveFile", saveFile, true, false)
	pages.AddPage("getFile", getFile, true, false)
	if err := app.SetRoot(pages, true).SetFocus(mainMenuPage).Run(); err != nil {
		panic(err)
	}
}
