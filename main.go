package main

import (
	"database/sql"
	"fmt"

	"github.com/ccb012100/go-playlist-search/config"
	"github.com/gdamore/tcell/v2"
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
	// Grid at Application Root
	grid *tview.Grid
	// pages shown with the application
	pages *tview.Pages
	// footer for displaying messages
	messageBar *tview.TextView
	// header for menu options
	menuBar *tview.TextView
	// db file path
	db string
}

func main() {
	conf := SetConfig()

	// create main View
	view := &View{
		db:         conf.DBFilePath,
		app:        tview.NewApplication().EnableMouse(true),
		grid:       tview.NewGrid(),
		pages:      tview.NewPages(),
		messageBar: tview.NewTextView(),
		menuBar:    tview.NewTextView(),
	}

	CreateRootGrid(view)
	CreatePages(view)

	view.UpdateMessageBar("Application created!")

	if err := view.app.SetRoot(view.grid, true).SetFocus(view.grid).Run(); err != nil {
		panic(err)
	}
}

func (v *View) UpdateMessageBar(message string) {
	v.messageBar.Clear().SetText(message)
}

func CreatePages(v *View) {
	CreateMenuBar(v)
	CreateMessageBar(v)
	CreateHomePage(v)
	CreatePlaylistsPage(v)
	CreateAlbumsPage(v)
	CreateAlbumsPage(v)
	CreateSongsPage(v)
}

func CreateMessageBar(v *View) {
	v.messageBar.SetTextAlign(tview.AlignCenter).SetText("Message Bar")
}

func CreateMenuBar(v *View) {
	v.menuBar.SetTextAlign(tview.AlignCenter).SetText("Menu Bar")
}

func CreateRootGrid(v *View) {
	// 3 rows
	// row 1: Menu Bar
	// row 2: main content
	// row 3: Message Bar
	v.grid.SetRows(3, 0, 5).SetColumns(0).SetBorders(true).
		// row 0
		AddItem(v.menuBar, 0, 0, 1, 1, 0, 0, false).
		// row 1
		AddItem(v.pages, 1, 0, 1, 1, 0, 0, true).
		// row 2
		AddItem(v.messageBar, 2, 0, 1, 1, 0, 0, false)
}

// Read configuration file and map it to a Config struct
func SetConfig() config.Config {
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
func IntToAlpha(i int) rune {
	return rune('a' - 1 + i)
}

func AddQuitOption(list *tview.List, f func()) {
	list.AddItem("[red::b]Quit[-]", "[red::]Press q to exit[-]", 'q', f)
}

func AddResetOption(list *tview.List, f func()) {
	list.AddItem("[yellow::b]Reset[-]", "[yellow::]Press r to reset this page[-]", 'r', f)
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
	input := tview.NewInputField()

	input.SetLabel("Search for a playlist: ").SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		v.UpdateMessageBar(fmt.Sprintf("key = %v", key))
		switch key {
		case tcell.KeyEscape:
			v.pages.ShowPage(HomePage)
			v.app.SetFocus(v.pages)
		case tcell.KeyEnter:
			ShowPlaylists(v, input.GetText())
		}
	})

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid := tview.NewGrid().SetRows(0, 0).SetBorders(true)
	grid.SetTitle("PLaylists")
	grid.AddItem(input, 0, 0, 1, 1, 0, 0, true)

	v.pages.AddPage(PlaylistsPage, grid, true, false)
}

func ShowPlaylists(v *View, query string) {
	v.UpdateMessageBar("show playlists")
	pageName, item := v.pages.GetFrontPage()

	if pageName != PlaylistsPage {
		panic("This method should only be called from the Playlists page")
	}

	grid, ok := item.(*tview.Grid)
	if !ok {
		panic(fmt.Sprintf("Expected type *tview.Grid but got %T", grid))
	}

	list := tview.NewList()
	database, _ := sql.Open("sqlite3", v.db)
	rows, err := database.Query(
		"SELECT id, name FROM Playlist WHERE name LIKE '%' || @Query || '%' ORDER BY name",
		sql.Named("Query", query))

	if err != nil {
		panic(err)
	}

	var i = 1
	for rows.Next() {
		var id string
		var name string

		rows.Scan(&id, &name)
		list.AddItem(name, id, 0, func() { SelectPlaylist(v, id, name) })
		i++
	}

	AddQuitOption(list, func() { v.pages.SwitchToPage(HomePage) })

	AddResetOption(list, func() {
		CreatePlaylistsPage(v)
		v.pages.ShowPage(PlaylistsPage)
		v.app.SetFocus(v.pages)
	})

	list.SetTitle("Playlists")

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid.Clear().AddItem(list, 0, 0, 1, 1, 0, 0, true)
	v.app.SetFocus(grid)
}

func SelectPlaylist(v *View, id string, name string) {
	v.UpdateMessageBar(fmt.Sprintf("Selected playlist %s %s", id, name))
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
