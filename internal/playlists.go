package pages

import (
	"database/sql"
	"fmt"

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
			ShowPlaylists(v, input.GetText())
		}
	})

	v.SetMainPanel(input)
}

func ShowPlaylists(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowPlaylists() query='%s'", query))

	database, _ := sql.Open("sqlite3", v.DB)
	rows, err := database.Query(
		"SELECT id, name FROM Playlist WHERE name LIKE '%' || @Query || '%' ORDER BY name",
		sql.Named("Query", query))

	if err != nil {
		panic(err)
	}

	var playlists []models.SimpleIdentifier

	for rows.Next() {
		var id string
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}

		playlists = append(playlists, models.SimpleIdentifier{Id: id, Name: name})
	}

	v.List.Clear().SetTitle("Playlists")

	// TODO: display message if there are no matches

	for _, plist := range playlists {
		p := plist
		v.List.AddItem(p.Name, p.Id, 0, func() { SelectPlaylist(v, p) })
	}

	AddQuitToHomeOption(v.List, v)
	AddResetOption(v.List, func() { SearchForPlaylists(v) })

	v.SetMainPanel(v.List)

	if len(playlists) == 1 {
		SelectPlaylist(v, playlists[0])
	}
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
			ShowStarredPlaylistMatches(v, input.GetText())
		}
	})

	v.SetMainPanel(input)
}

func ShowStarredPlaylistMatches(v *models.View, query string) {
	v.UpdateTitleBar(fmt.Sprintf("Items in Starred Playlists matching '%s'", query))
	database, _ := sql.Open("sqlite3", v.DB)

	// Get albums by the artist
	albumArtistRows, err := database.Query(
		"SELECT P.name AS playlistName, T.name AS trackName, A.name AS albumName, GROUP_CONCAT(A2.name, '; ') AS artists FROM Playlist P JOIN PlaylistTrack PT ON P.id = PT.playlist_id JOIN Track T ON PT.track_id = T.id JOIN Album A ON T.album_id = A.id JOIN TrackArtist TA ON T.id = TA.track_id JOIN Artist A2 ON TA.artist_id = A2.id WHERE P.name LIKE 'Starred%' AND (A2.name LIKE '%' || @Query || '%' OR T.name LIKE '%' || @Query || '%' OR A.name LIKE '%' || @Query || '%') GROUP BY P.name, T.id, A.id, PT.added_at, T.track_number ORDER BY P.name, A.id, PT.added_at, T.track_number",
		sql.Named("Query", query))

	if err != nil {
		panic(err)
	}

	var matches []models.StarredPlaylistMatch

	for albumArtistRows.Next() {
		var playlistName, trackName, albumName, artists string

		if err := albumArtistRows.Scan(&playlistName, &trackName, &albumName, &artists); err != nil {
			panic(err)
		}

		matches = append(matches, models.StarredPlaylistMatch{
			PlaylistName: playlistName,
			TrackName:    trackName,
			AlbumName:    albumName,
			Artists:      artists,
		})
	}

	if len(matches) == 0 {
		textView := tview.NewTextView().SetDynamicColors(true)
		textView.SetTitle("No matches").SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
		textView.SetText(fmt.Sprintf("There are no matches for the query [green:-:b]%s[-]", query))

		textView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
			switch e.Key() {
			case tcell.KeyESC:
				v.SetMainPanel(v.List)
				// case tcell.KeyRune:
				// 	v.SetMainPanel(v.List)
			}

			return e
		})

		v.SetMainPanel(textView)
		return
	}

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
		table.SetCell(i+1, 1, tview.NewTableCell(padLeft(match.TrackName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignRight).SetExpansion(1))
		table.SetCell(i+1, 2, tview.NewTableCell(padLeft(match.AlbumName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 3, tview.NewTableCell(padLeft(match.Artists)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	table.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(table)
}
