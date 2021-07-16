package main

import (
	pages "github.com/ccb012100/go-playlist-search/internal"
	"github.com/ccb012100/go-playlist-search/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

// TODO: add key shortcuts to switch focus between grid sections
func main() {
	conf := pages.SetConfig()

	// create main View
	view := &models.View{
		DB:         conf.DBFilePath,
		App:        tview.NewApplication().EnableMouse(true),
		Grid:       tview.NewGrid(),
		Pages:      tview.NewPages(),
		MessageBar: tview.NewTextView(),
		MenuBar:    tview.NewTextView(),
	}

	pages.CreateRootGrid(view)
	pages.CreatePages(view)

	view.UpdateMessageBar("Application created!")

	if err := view.App.SetRoot(view.Grid, true).SetFocus(view.Grid).Run(); err != nil {
		panic(err)
	}
}
