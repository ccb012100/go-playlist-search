package pages

import (
	"database/sql"
	"fmt"
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

	database, _ := sql.Open("sqlite3", v.DB)
	// NOTE: this query only returns matches from Album Artists, not Track Artists
	sqlRows, err := database.Query(
		"select id, name, total_tracks, release_date, album_type from Album a join AlbumArtist AA on a.id = AA.album_id where AA.artist_id = @Id order by release_date",
		sql.Named("Id", artist.Id))

	if err != nil {
		panic(err)
	}

	var albums []models.Album

	for sqlRows.Next() {
		var id, name, releaseDate, albumType string
		var totalTracks int

		if err := sqlRows.Scan(&id, &name, &totalTracks, &releaseDate, &albumType); err != nil {
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

	table := tview.NewTable().SetBorders(true)
	// set header row
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 1, tview.NewTableCell("Tracks").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 2, tview.NewTableCell("Release Date").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 3, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))

	// TODO: add left/right padding to table cells
	// set table contents
	for r := 0; r < len(albums); r++ {
		album := &albums[r]
		// use r+1 to offset for header row
		table.SetCell(r+1, 0, tview.NewTableCell(album.Name).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(r+1, 1, tview.NewTableCell(strconv.Itoa(album.TotalTracks)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignRight).SetExpansion(1))
		table.SetCell(r+1, 2, tview.NewTableCell(album.ReleaseDate).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(r+1, 3, tview.NewTableCell(album.AlbumType).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	table.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(table)
}
