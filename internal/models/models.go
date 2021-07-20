package models

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

type View struct {
	// main application
	App *tview.Application
	// Grid at Application Root
	Grid *tview.Grid
	// footer for displaying messages
	MessageBar *tview.TextView
	// header at top of app
	TitleBar *tview.TextView
	// db file path
	DB string
	// Selection List
	List *tview.List
}

type Album struct {
	Id          string
	Name        string
	TotalTracks int
	ReleaseDate string
	AlbumType   string
}

type SimpleIdentifier struct {
	Name string
	Id   string
}

func (v View) UpdateMessageBar(message string) {
	v.MessageBar.SetText(fmt.Sprintf("%s => %s", time.Now().Format("03:04:05"), message))
}

func (v View) UpdateTitleBar(message string) {
	v.TitleBar.SetText(message)
}

// Display the Primitive in the main panel of the app's Grid
func (v View) SetMainPanel(p tview.Primitive) {
	v.Grid.AddItem(p, 1, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(p)
}
