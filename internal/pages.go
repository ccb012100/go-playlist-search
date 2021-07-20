package pages

import (
	"fmt"

	"github.com/ccb012100/go-playlist-search/internal/models"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

const (
	ALBUMS_PAGE    = "ALBUMS_PAGE"
	ARTISTS_PAGE   = "ARTISTS_PAGE"
	HOME_PAGE      = "HOME_PAGE"
	PLAYLISTS_PAGE = "PLAYLISTS_PAGE"
	SONGS_PAGE     = "SONGS_PAGE"
)

func CreatePages(v *models.View) {
	CreateMenuBar(v)
	CreateMessageBar(v)
	CreateHomePage(v)
	CreatePlaylistsPage(v)
	CreateArtistsPage(v)
	CreateAlbumsPage(v)
	CreateSongsPage(v)
}

func CreateMessageBar(v *models.View) {
	v.MessageBar.SetTextAlign(tview.AlignCenter).SetText("Message Bar")
}

func CreateMenuBar(v *models.View) {
	v.MenuBar.SetTextAlign(tview.AlignCenter).SetText("Menu Bar")
}

func CreateViewGrid(v *models.View) {
	v.Grid.SetRows(2, 0, 2).SetColumns(0).SetBorders(true).
		// row 0: Menu Bar
		AddItem(v.MenuBar, 0, 0, 1, 1, 0, 0, false).
		// row 1: main content
		AddItem(v.Pages, 1, 0, 1, 1, 0, 0, true).
		// row 2: Message Bar
		AddItem(v.MessageBar, 2, 0, 1, 1, 0, 0, false)
}

func CreateHomePage(v *models.View) {
	list := tview.NewList().
		AddItem("Playlists", "View Playlists", '1', func() { v.Pages.SwitchToPage(PLAYLISTS_PAGE) }).
		AddItem("Artists", "View Artists", '2', func() { v.Pages.SwitchToPage(ARTISTS_PAGE) }).
		AddItem("Albums", "View Albums", '3', func() { v.Pages.SwitchToPage(ALBUMS_PAGE) }).
		AddItem("Songs", "View Songs", '4', func() { v.Pages.SwitchToPage(SONGS_PAGE) })

	AddQuitOption(list, func() { v.App.Stop() })

	list.SetTitle("Home")

	v.Pages.AddPage(HOME_PAGE, list, true, true)
}

func CreatePlaylistsPage(v *models.View) {
	input := tview.NewInputField()

	input.SetLabel("Search for a playlist: ").SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		v.UpdateMessageBar(fmt.Sprintf("key = %v", key))
		switch key {
		case tcell.KeyEscape:
			v.Pages.ShowPage(HOME_PAGE)
			v.App.SetFocus(v.Pages)
		case tcell.KeyEnter:
			ShowPlaylists(v, input.GetText())
		}
	})

	v.Pages.AddPage(PLAYLISTS_PAGE, input, true, false)
}

func CreateArtistsPage(v *models.View) {
	input := tview.NewInputField()

	input.SetLabel("Search for artists: ").SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		v.UpdateMessageBar(fmt.Sprintf("key = %v", key))
		switch key {
		case tcell.KeyEscape:
			v.Pages.ShowPage(HOME_PAGE)
			v.App.SetFocus(v.Pages)
		case tcell.KeyEnter:
			ShowArtists(v, input.GetText())
		}
	})

	v.Pages.AddPage(ARTISTS_PAGE, input, true, false)
}

func CreateSongsPage(v *models.View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Songs")
	v.Pages.AddPage(SONGS_PAGE, box, true, false)
}

func CreateAlbumsPage(v *models.View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Albums")
	v.Pages.AddPage(SONGS_PAGE, box, true, false)
}

// convert ints into alphabetic characters,
// i.e. 1->a, 2->b, 3->c, etc.
func IntToAlpha(i int) rune {
	return rune('a' - 1 + i)
}

// add a Quit option to the passed-in list
func AddQuitOption(list *tview.List, f func()) {
	list.AddItem("[red::b]Quit[-]", "[red::]Press q to exit[-]", 'q', f)
}

// add a Quit to Home Page option to the passed-in list
func AddQuitToHomeOption(list *tview.List, v *models.View) {
	AddQuitOption(list, func() {
		v.Pages.SwitchToPage(HOME_PAGE)
		v.App.SetFocus(v.Pages)
	})
}

// add a Reset Page option to the passed-in list
func AddResetOption(list *tview.List, f func()) {
	list.AddItem("[yellow::b]Reset[-]", "[yellow::]Press r to reset this page[-]", 'r', f)
}
