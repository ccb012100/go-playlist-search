package main

import (
	"github.com/ccb012100/go-playlist-search/config"
	"github.com/ccb012100/go-playlist-search/internal"
	"github.com/ccb012100/go-playlist-search/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

func main() {
	conf := config.SetConfig()

	// create main View
	view := &models.View{
		DB:  conf.DBFilePath,
		App: tview.NewApplication().EnableMouse(true),
	}

	internal.CreateViewGrid(view)
	internal.GoToMainMenu(view)

	view.UpdateMessageBar("Application created!")

	if err := view.App.SetRoot(view.Grid, true).SetFocus(view.Grid).Run(); err != nil {
		panic(err)
	}
}
