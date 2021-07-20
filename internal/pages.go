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

// Add key bindings for selecting list items.
func AddListInputListener(l *tview.List) {
	l.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyRune:
			switch e.Rune() {
			case rune('n'):
				SelectNextListItem(l)
			case 'p':
				SelectPreviousListIten(l)
			case 'j':
				SelectNextListItem(l)
			case 'k':
				SelectPreviousListIten(l)
			}
		}

		return e
	})
}

// Select Next List Item.
// Wrap around to the first item if the last is currently selected.
func SelectNextListItem(l *tview.List) {
	i := l.GetCurrentItem() + 1

	if i < l.GetItemCount() {
		l.SetCurrentItem(i)
	} else {
		l.SetCurrentItem(0)
	}
}

// Select Previous List Item.
// Wrap around to last item if the first item is currently selected.
func SelectPreviousListIten(l *tview.List) {
	l.SetCurrentItem(l.GetCurrentItem() - 1)
}

// Display the Primitive in the main panel of the app's Grid
func SetMainPanel(p tview.Primitive, v *models.View) {
	v.Grid.AddItem(p, 1, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(p)
	v.UpdateMessageBar(fmt.Sprintf("p has focus => %t", p.HasFocus()))
}

func CreateMessageBar(v *models.View) {
	v.MessageBar = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Message Bar")
	v.MessageBar.SetBorder(true).SetBorderColor(tcell.ColorDarkGreen)
}

func CreateMenuBar(v *models.View) {
	v.MenuBar = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Menu Bar")
	v.MenuBar.SetBorder(true).SetBorderColor(tcell.ColorHotPink)
}

func CreateViewGrid(v *models.View) {
	CreateMenuBar(v)
	CreateMessageBar(v)

	v.List = tview.NewList()
	v.List.SetBorder(true).SetBorderColor(tcell.ColorDarkRed).SetTitle("List")
	AddListInputListener(v.List)

	v.Grid = tview.NewGrid().SetRows(2, 0, 2).SetColumns(0).SetBorders(true)
	v.Grid.SetBorderColor(tcell.ColorMediumPurple)

	// row 0: Menu Bar
	v.Grid.AddItem(v.MenuBar, 0, 0, 1, 1, 0, 0, false).
		// row 1: main content
		AddItem(v.List, 1, 0, 1, 1, 0, 0, true).
		// row 2: Message Bar
		AddItem(v.MessageBar, 2, 0, 1, 1, 0, 0, false)
}

func GoToMainMenu(v *models.View) {
	v.List.Clear().
		AddItem("Playlists", "View Playlists", '1', func() {}).
		AddItem("Artists", "View Artists", '2', func() {}).
		AddItem("Albums", "View Albums", '3', func() {}).
		AddItem("Songs", "View Songs", '4', func() {})

	AddQuitOption(v.List, func() { v.App.Stop() })

	v.List.SetTitle("Main Menu").SetBorderColor(tcell.ColorDarkRed)

	SetMainPanel(v.List, v)
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
