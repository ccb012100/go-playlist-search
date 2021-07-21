package pages

import (
	"fmt"
	"strconv"

	"github.com/ccb012100/go-playlist-search/internal/data"
	"github.com/ccb012100/go-playlist-search/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SearchForArtists(v *models.View) {
	input := tview.NewInputField()
	// TODO: set minimum input length
	input.SetLabel("Search for artists: ").SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		v.UpdateMessageBar(fmt.Sprintf("key = %v", key))
		switch key {
		case tcell.KeyEscape:
			GoToMainMenu(v)
		case tcell.KeyEnter:
			ShowArtistSearchResults(v, input.GetText())
		}
	})

	v.SetMainPanel(input)
}

func ShowArtistSearchResults(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowArtists() query='%s'", query))
	v.UpdateTitleBar(fmt.Sprintf("Artists matching '%s'", query))

	artists := data.SearchArtists(query, v.DB)

	// show message if 0 results
	if len(artists) == 0 {
		textView := tview.NewTextView().SetDynamicColors(true)
		textView.SetTitle("No matches").SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
		textView.SetText(fmt.Sprintf("There are no Artists matching [green:-:b]%s[-]", query))

		textView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
			switch e.Key() {
			case tcell.KeyESC:
				v.SetMainPanel(v.List)
			}

			return e
		})

		v.SetMainPanel(textView)
		return
	}

	// if there's only 1 match, just select it
	if len(artists) == 1 {
		SelectArtist(v, artists[0])
		return
	}

	displayArtists(v, artists)
}

func displayArtists(v *models.View, artists []models.SimpleIdentifier) {
	v.List.Clear()
	for _, artist := range artists {
		a := artist
		v.List.AddItem(artist.Name, artist.Id, 0, func() { SelectArtist(v, a) })
	}
	AddQuitToHomeOption(v.List, v)

	AddResetOption(v.List, func() {
		SearchForArtists(v)
	})

	v.List.SetTitle("Artists results")

	v.SetMainPanel(v.List)
}

func SelectArtist(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateTitleBar(artist.Name)
	v.UpdateMessageBar(fmt.Sprintf("Selected artist %s %s", artist.Id, artist.Name))

	v.List.Clear().
		AddItem("Albums", "View Artist's Albums", '1', func() { ShowArtistAlbums(v, artist) }).
		AddItem("Tracks", "View Artist's Tracks", '2', func() { v.UpdateMessageBar("Not implemented yet") }).
		AddItem("Playlists", "List Playlists containing the Artist", '3', func() { showPlaylistsWithArtist(v, artist) })

	AddQuitOption(v.List, func() { GoToMainMenu(v) })

	v.List.SetTitle("Artist Info").SetBorderColor(tcell.ColorDarkSeaGreen)

	v.SetMainPanel(v.List)
}

func ShowArtistAlbums(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateTitleBar(fmt.Sprintf("Albums by %s", artist.Name))

	albums := data.GetAlbumsByArtist(&artist, v.DB)

	// Display message if there are no albums found
	if len(albums) == 0 {
		textView := tview.NewTextView().SetDynamicColors(true)
		textView.SetTitle("No matches").SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
		textView.SetText(fmt.Sprintf("There are no Albums for artist [green:-:b]%s[-] [gray:-:-](Id = %s)[-]", artist.Name, artist.Id))

		textView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
			switch e.Key() {
			case tcell.KeyESC:
				v.SetMainPanel(v.List)
			}

			return e
		})

		v.SetMainPanel(textView)
		return
	}

	displayArtistAlbumsTable(v, albums)
}

func displayArtistAlbumsTable(v *models.View, albums []models.Album) {
	table := tview.NewTable().SetBorders(true)
	// set header row
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 1, tview.NewTableCell("Tracks").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 2, tview.NewTableCell("Release Date").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 3, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))

	// set table contents
	for i := 0; i < len(albums); i++ {
		album := albums[i]

		// use i+1 to offset for header
		table.SetCell(i+1, 0, tview.NewTableCell(padLeft(album.Name)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(i+1, 1, tview.NewTableCell(padRight(strconv.Itoa(album.TotalTracks))).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignRight).SetExpansion(1))
		table.SetCell(i+1, 2, tview.NewTableCell(padLeft(album.ReleaseDate)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(i+1, 3, tview.NewTableCell(padLeft(album.AlbumType)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	table.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(table)
}

// Display Playlists containing the specified Artist
func showPlaylistsWithArtist(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateTitleBar("Playlists containing tracks by " + artist.Name)

	playlists := data.FindPlaylistsContainingArtist(artist, v.DB)

	// this should never happen
	if len(playlists) == 0 {
		panic(fmt.Sprintf("No playlists were found for artist '%s', '%s'", artist.Name, artist.Id))
	}

	txt := fmt.Sprintf("Playlists containing %s:\n", artist.Name)

	for i, p := range playlists {
		txt += fmt.Sprintf("\n%d.\t%s\t%s", i+1, p.Id, p.Name)
	}

	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetText(txt)
	textView.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(textView)
}
