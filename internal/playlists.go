package internal

import (
	"fmt"

	"github.com/ccb012100/go-playlist-search/internal/data"
	"github.com/ccb012100/go-playlist-search/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SearchForPlaylists(v *models.View) {
	input := tview.NewInputField()
	// TODO: set minimum input length
	input.SetLabel("Search for a playlist: ").SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		v.UpdateMessageBar(fmt.Sprintf("key = %v", key))
		switch key {
		case tcell.KeyEscape:
			GoToMainMenu(v)
		case tcell.KeyEnter:
			ShowPlaylistSearchResults(v, input.GetText())
		}
	})

	v.SetMainPanel(input)
}

func ShowPlaylistSearchResults(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowPlaylists() query='%s'", query))

	playlists := data.SearchPlaylists(query, v.DB)

	// display message if there are no matches
	if len(playlists) == 0 {
		textView := tview.NewTextView().SetDynamicColors(true)
		textView.SetTitle("No matches").SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
		textView.SetText(fmt.Sprintf("There are no matches for the query [green:-:b]%s[-]", query))

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

	displayPlaylists(v, playlists)
}

func displayPlaylists(v *models.View, playlists []models.SimpleIdentifier) {
	// if there's only 1 match, just select it
	if len(playlists) == 1 {
		SelectPlaylist(v, playlists[0])
		return
	}

	v.List.Clear().SetTitle("Playlists")
	for _, plist := range playlists {
		p := plist
		v.List.AddItem(p.Name, p.Id, 0, func() { SelectPlaylist(v, p) })
	}
	AddQuitToHomeOption(v.List, v)
	AddResetOption(v.List, func() { SearchForPlaylists(v) })

	v.SetMainPanel(v.List)
}

func SelectPlaylist(v *models.View, playlist models.SimpleIdentifier) {
	v.UpdateTitleBar(playlist.Name)
	var textView = tview.NewTextView()

	textView.SetText(fmt.Sprintf("Selected playlist id='%s', name='%s'", playlist.Id, playlist.Name)).
		SetTitle(fmt.Sprintf("Selected Playlist: %s", playlist.Name))

	textView.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(textView)
}

func SearchStarredPlaylists(v *models.View) {
	input := tview.NewInputField()
	// TODO: set minimum input length
	input.SetLabel("Search Starred playlists: ").SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		v.UpdateMessageBar(fmt.Sprintf("key = %v", key))
		switch key {
		case tcell.KeyEscape:
			GoToMainMenu(v)
		case tcell.KeyEnter:
			ShowStarredPlaylistSearchResults(v, input.GetText())
		}
	})

	v.SetMainPanel(input)
}

func ShowStarredPlaylistSearchResults(v *models.View, query string) {
	v.UpdateTitleBar(fmt.Sprintf("Items in Starred Playlists matching '%s'", query))

	matches := data.SearchStarredPlaylists(query, v.DB)

	if len(matches) == 0 {
		textView := tview.NewTextView().SetDynamicColors(true)
		textView.SetTitle("No matches").SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
		textView.SetText(fmt.Sprintf("There are no matches for the query [green:-:b]%s[-]", query))

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

	displayStarredPlaylistMatches(v, matches)
}

func displayStarredPlaylistMatches(v *models.View, matches []models.StarredPlaylistMatch) {
	table := tview.NewTable().SetBorders(true)
	// set header row
	table.SetCell(0, 0, tview.NewTableCell("Playlist").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 1, tview.NewTableCell("Track").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 2, tview.NewTableCell("Album").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 3, tview.NewTableCell("Artists").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))

	// set table contents
	// start row at 1 to offset for header
	for i := 0; i < len(matches); i++ {
		match := matches[i]

		// use i+1 to offset for header row
		table.SetCell(i+1, 0, tview.NewTableCell(padLeft(match.PlaylistName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 1, tview.NewTableCell(padLeft(match.TrackName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 2, tview.NewTableCell(padLeft(match.AlbumName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 3, tview.NewTableCell(padLeft(match.Artists)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	table.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(table)
}
