package main

import (
	"database/sql"
	"fmt"

	"github.com/ccb012100/go-playlist-search/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

const (
	AlbumsPage    = "albums-page"
	ArtistsPage   = "artists-page"
	HomePage      = "home-page"
	PlaylistsPage = "playlists-page"
	SongsPage     = "songs-page"
)

type View struct {
	// main application
	app *tview.Application
	// pages shown with the application
	pages *tview.Pages
	// panel at bottom of app for displaying messages
	messageBar *tview.TextView
	// db file path
	db string
}

func main() {
	conf := setConfig()

	// create main View
	view := &View{
		db:         conf.DBFilePath,
		app:        tview.NewApplication().EnableMouse(true),
		pages:      tview.NewPages(),
		messageBar: tview.NewTextView(),
	}

	CreateHomePage(view)
	CreatePlaylistsPage(view)
	CreateAlbumsPage(view)
	CreateAlbumsPage(view)
	CreateSongsPage(view)

	// TODO: create Grid: Pages in Top Row, MessageBar in bottom row for displaying log messages

	if err := view.app.SetRoot(view.pages, true).SetFocus(view.pages).Run(); err != nil {
		panic(err)
	}
}

// Read configuration file and map it to a Config struct
func setConfig() config.Config {
	viper.New()
	viper.SetConfigFile("./app.env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var configuration config.Config

	if err := viper.Unmarshal(&configuration); err != nil {
		panic(err)
	}

	return configuration
}

// convert ints into alphabetic characters,
// i.e. 1->a, 2->b, 3->c, etc.
func intToAlpha(i int) rune {
	return rune('a' - 1 + i)
}

func AddQuitOption(list *tview.List, f func()) {
	list.AddItem("Quit", "[red]Press [::b]q[::-] to exit[-]", 'q', f)
}

func CreateHomePage(v *View) {
	list := tview.NewList().
		AddItem("Playlists", "View Playlists", '1', func() { v.pages.SwitchToPage(PlaylistsPage) }).
		AddItem("Artists", "View Artists", '2', func() { v.pages.SwitchToPage(ArtistsPage) }).
		AddItem("Albums", "View Albums", '3', func() { v.pages.SwitchToPage(AlbumsPage) }).
		AddItem("Songs", "View Songs", '4', func() { v.pages.SwitchToPage(SongsPage) })

	AddQuitOption(list, func() { v.app.Stop() })

	list.SetTitle("Home")

	v.pages.AddPage(HomePage, list, true, true)
}

func CreatePlaylistsPage(v *View) {
	list := tview.NewList()
	database, _ := sql.Open("sqlite3", v.db)
	rows, _ := database.Query("SELECT id, name FROM Playlist LIMIT 10")

	var i = 1
	for rows.Next() {
		var id string
		var name string

		rows.Scan(&id, &name)
		list.AddItem(name, id, intToAlpha(i), func() { SelectPlaylist(v, id, name) })
		i++
	}
	AddQuitOption(list, func() { v.pages.SwitchToPage(HomePage) })

	list.SetTitle("Playlists")

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid := tview.NewGrid().SetRows(0, 0).SetBorders(true)
	grid.AddItem(list, 0, 0, 1, 1, 0, 0, true)

	v.pages.AddPage(PlaylistsPage, grid, true, false)
}

func SelectPlaylist(v *View, id string, name string) {
	pageName, item := v.pages.GetFrontPage()

	if pageName != PlaylistsPage {
		panic("This method should only be called from the Playlists page")
	}

	grid, ok := item.(*tview.Grid)
	if !ok {
		panic(fmt.Sprintf("Expected type *tview.Grid but got %T", grid))
	}

	var textView = tview.NewTextView()

	textView.SetText(fmt.Sprintf("Selected playlist id='%s', name='%s'", id, name)).
		SetTitle(fmt.Sprintf("Selected Playlist: %s", name))

	grid.AddItem(textView, 1, 0, 1, 1, 0, 0, false)
}

func CreateArtistsPage(v *View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Artists")
	v.pages.AddPage(ArtistsPage, box, true, false)
}

func CreateSongsPage(v *View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Songs")
	v.pages.AddPage(SongsPage, box, true, false)
}

func CreateAlbumsPage(v *View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Albums")
	v.pages.AddPage(SongsPage, box, true, false)
}

func SelectArtist(v *View, id string, name string) {

}

func SelectSong(v *View, id string, name string) {

}

func SelectAlbum(v *View, id string, name string) {

}
