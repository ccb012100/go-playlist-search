package pages

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ccb012100/go-playlist-search/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ShowArtists(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowArtists() query='%s'", query))
	pageName, _ := v.Pages.GetFrontPage()

	if pageName != ARTISTS_PAGE {
		panic("This method should only be called from the Playlists page")
	}

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

	list := tview.NewList()
	for _, artist := range artists {
		a := artist
		list.AddItem(artist.Name, artist.Id, 0, func() { SelectArtist(v, a) })
	}
	// TODO: show message if 0 results
	AddQuitToHomeOption(list, v)

	AddResetOption(list, func() {
		CreateArtistsPage(v)
		v.Pages.ShowPage(ARTISTS_PAGE)
		v.App.SetFocus(v.Pages)
	})

	list.SetTitle("Artists results")

	v.Grid.AddItem(list, 1, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(list)

	if len(artists) == 1 {
		SelectArtist(v, artists[0])
	}
}

func SelectArtist(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateMessageBar(fmt.Sprintf("Selected artist %s %s", artist.Id, artist.Name))
	pageName, _ := v.Pages.GetFrontPage()

	if pageName != ARTISTS_PAGE {
		panic(fmt.Sprintf("This method should only be called from %s, but was called from %s", ARTISTS_PAGE, pageName))
	}

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

		v.Grid.AddItem(textView, 1, 0, 1, 1, 0, 0, true)
		v.App.SetFocus(textView)

		// TODO: add key binding to go back to select ("press any key to go back")
		return
	}

	table := tview.NewTable().SetBorders(true)
	// set header row
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 1, tview.NewTableCell("Tracks").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 2, tview.NewTableCell("Release Date").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(2))
	table.SetCell(0, 3, tview.NewTableCell("Type").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))

	// set table contents
	for r := 0; r < len(albums); r++ {
		album := &albums[r]
		// use r+1 to offset for header row
		table.SetCell(r+1, 0, tview.NewTableCell(album.Name).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(r+1, 1, tview.NewTableCell(strconv.Itoa(album.TotalTracks)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignRight).SetExpansion(1))
		table.SetCell(r+1, 2, tview.NewTableCell(album.ReleaseDate).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(2))
		table.SetCell(r+1, 3, tview.NewTableCell(album.AlbumType).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	v.Grid.AddItem(table, 1, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(table)
}
