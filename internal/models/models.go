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
	// pages shown with the application
	Pages *tview.Pages
	// footer for displaying messages
	MessageBar *tview.TextView
	// header for menu options
	MenuBar *tview.TextView
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
	v.MessageBar.Clear().SetText(fmt.Sprintf("%s => %s", time.Now().Format("03:04:05"), message))
}
