package pages

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

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
			ShowArtists(v, input.GetText())
		}
	})

	v.SetMainPanel(input)
}

func ShowArtists(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowArtists() query='%s'", query))
	v.UpdateTitleBar(fmt.Sprintf("Artists matching '%s'", query))

	database, _ := sql.Open("sqlite3", v.DB)

	rows, err := database.Query(
		"SELECT id, name FROM Artist WHERE name LIKE '%' || @Query || '%' ORDER BY name",
		sql.Named("Query", query))

	if err != nil {
		panic(err)
	}

	var artists []models.SimpleIdentifier

	for rows.Next() {
		var id string
		var name string

		rows.Scan(&id, &name)
		artists = append(artists, models.SimpleIdentifier{Id: id, Name: name})
	}

	// TODO: show message if 0 results
	if len(artists) == 1 {
		SelectArtist(v, artists[0])
		return
	}

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
		AddItem("Playlists", "List Playlists containing the Artist", '3', func() { showPlaylists(v, artist) })

	AddQuitOption(v.List, func() { GoToMainMenu(v) })

	v.List.SetTitle("Artist Info").SetBorderColor(tcell.ColorDarkSeaGreen)

	v.SetMainPanel(v.List)
}

func ShowArtistAlbums(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateTitleBar(fmt.Sprintf("Albums by %s", artist.Name))
	database, _ := sql.Open("sqlite3", v.DB)

	// Get albums by the artist
	albumArtistRows, err := database.Query(
		"select id, name, total_tracks, release_date, album_type from Album a join AlbumArtist AA on a.id = AA.album_id where AA.artist_id = @Id",
		sql.Named("Id", artist.Id))

	if err != nil {
		panic(err)
	}

	var albums []models.Album

	for albumArtistRows.Next() {
		var id, name, releaseDate, albumType string
		var totalTracks int

		if err := albumArtistRows.Scan(&id, &name, &totalTracks, &releaseDate, &albumType); err != nil {
			panic(err)
		}

		albums = append(albums, models.Album{
			Id:          id,
			Name:        name,
			TotalTracks: totalTracks,
			ReleaseDate: releaseDate,
			AlbumType:   albumType,
		})
	}

	// Get albums with Tracks the Artist appears on
	trackArtistRows, err := database.Query(
		"select A.id, A.name, total_tracks, release_date, album_type from Album A join Track T on A.id = T.album_id join TrackArtist TA on T.id = TA.track_id where TA.artist_id = @Id",
		sql.Named("Id", artist.Id))

	if err != nil {
		panic(err)
	}

	for trackArtistRows.Next() {
		var id, name, releaseDate, albumType string
		var totalTracks int

		if err := trackArtistRows.Scan(&id, &name, &totalTracks, &releaseDate, &albumType); err != nil {
			panic(err)
		}

		albums = append(albums, models.Album{
			Id:          id,
			Name:        name,
			TotalTracks: totalTracks,
			ReleaseDate: releaseDate,
			AlbumType:   albumType,
		})
	}

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

	// sort albums
	sort.Sort(models.ByReleaseDate(albums))

	// track albums in a map so that we display a unique set
	var set = make(map[string]models.Album)

	table := tview.NewTable().SetBorders(true)
	// set header row
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 1, tview.NewTableCell("Tracks").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 2, tview.NewTableCell("Release Date").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 3, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))

	// set table contents
	// start row at 1 to offset for header
	for row, i := 1, 0; i < len(albums); i++ {
		album := albums[i]

		// skip if the album has already been added in
		if _, ok := set[album.Id]; ok {
			continue
		}

		set[album.Id] = album

		table.SetCell(row, 0, tview.NewTableCell(padLeft(album.Name)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(row, 1, tview.NewTableCell(padRight(strconv.Itoa(album.TotalTracks))).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignRight).SetExpansion(1))
		table.SetCell(row, 2, tview.NewTableCell(padLeft(album.ReleaseDate)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(row, 3, tview.NewTableCell(padLeft(album.AlbumType)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))

		row++
	}

	table.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(table)
}

// Display Playlists containing the specified Artist
func showPlaylists(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateTitleBar("Playlists containing tracks by " + artist.Name)
	database, _ := sql.Open("sqlite3", v.DB)
	sqlRows, err := database.Query(
		"select PL.name, PL.id from Playlist PL join PlaylistTrack PT on PL.id = PT.playlist_id join Track T on PT.track_id = T.id join TrackArtist TA on T.id = TA.track_id where TA.artist_id = @Id group by PL.id, PL.name order by Pl.name",
		sql.Named("Id", artist.Id))

	if err != nil {
		panic(err)
	}

	var playlists []models.SimpleIdentifier

	for sqlRows.Next() {
		var id, name string

		if err := sqlRows.Scan(&id, &name); err != nil {
			panic(err)
		}

		playlists = append(playlists, models.SimpleIdentifier{
			Id:   id,
			Name: name,
		})
	}

	textView := tview.NewTextView().SetDynamicColors(true)

	// this should never happen
	if len(playlists) == 0 {
		panic(fmt.Sprintf("No playlists were found for artist '%s', '%s'", artist.Name, artist.Id))
	}

	txt := fmt.Sprintf("Playlists containing %s:\n", artist.Name)

	for i, p := range playlists {
		txt += fmt.Sprintf("\n%d.\t%s\t%s", i, p.Id, p.Name)
	}

	textView.SetText(txt)

	textView.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(textView)
}
