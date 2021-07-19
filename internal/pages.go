package pages

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ccb012100/go-playlist-search/config"
	"github.com/ccb012100/go-playlist-search/internal/models"
	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
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

func CreateRootGrid(v *models.View) {
	v.Grid.SetRows(2, 0, 2).SetColumns(0).SetBorders(true).
		// row 0: Menu Bar
		AddItem(v.MenuBar, 0, 0, 1, 1, 0, 0, false).
		// row 1: main content
		AddItem(v.Pages, 1, 0, 1, 1, 0, 0, true).
		// row 2: Message Bar
		AddItem(v.MessageBar, 2, 0, 1, 1, 0, 0, false)
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

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid := tview.NewGrid().SetRows(0, 0).SetBorders(true)
	grid.SetTitle("PLaylists")
	grid.AddItem(input, 0, 0, 1, 1, 0, 0, true)

	v.Pages.AddPage(PLAYLISTS_PAGE, grid, true, false)
}

func ShowPlaylists(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowPlaylists() query='%s'", query))
	pageName, item := v.Pages.GetFrontPage()

	if pageName != PLAYLISTS_PAGE {
		panic(fmt.Sprintf("This method should only be called from %s, but was called from %s", PLAYLISTS_PAGE, pageName))
	}

	grid, ok := item.(*tview.Grid)
	if !ok {
		panic(fmt.Sprintf("Expected type *tview.Grid but got %T", grid))
	}

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

	list := tview.NewList()
	for _, plist := range playlists {
		list.AddItem(plist.Name, plist.Id, 0, func() { SelectPlaylist(v, plist) })
	}

	AddQuitOption(list, func() { v.Pages.SwitchToPage(HOME_PAGE) })

	AddResetOption(list, func() {
		CreatePlaylistsPage(v)
		v.Pages.ShowPage(PLAYLISTS_PAGE)
		v.App.SetFocus(v.Pages)
	})

	list.SetTitle("Playlists")

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid.Clear().AddItem(list, 0, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(grid)
}

func SelectPlaylist(v *models.View, playlist models.SimpleIdentifier) {
	pageName, item := v.Pages.GetFrontPage()

	if pageName != PLAYLISTS_PAGE {
		panic("This method should only be called from the Playlists page")
	}

	grid, ok := item.(*tview.Grid)
	if !ok {
		panic(fmt.Sprintf("Expected type *tview.Grid but got %T", grid))
	}

	var textView = tview.NewTextView()

	textView.SetText(fmt.Sprintf("Selected playlist id='%s', name='%s'", playlist.Id, playlist.Name)).
		SetTitle(fmt.Sprintf("Selected Playlist: %s", playlist.Name))

	grid.AddItem(textView, 1, 0, 1, 1, 0, 0, false)
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

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid := tview.NewGrid().SetRows(20, 30).SetBorders(true)
	grid.SetTitle("Artists")
	grid.AddItem(input, 0, 0, 1, 1, 0, 0, true)

	v.Pages.AddPage(ARTISTS_PAGE, grid, true, false)
}

func ShowArtists(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowArtists() query='%s'", query))
	pageName, item := v.Pages.GetFrontPage()

	if pageName != ARTISTS_PAGE {
		panic("This method should only be called from the Playlists page")
	}

	grid, ok := item.(*tview.Grid)
	if !ok {
		panic(fmt.Sprintf("Expected type *tview.Grid but got %T", grid))
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
	AddQuitOption(list, func() { v.Pages.SwitchToPage(HOME_PAGE) })

	AddResetOption(list, func() {
		CreateArtistsPage(v)
		v.Pages.ShowPage(ARTISTS_PAGE)
		v.App.SetFocus(v.Pages)
	})

	list.SetTitle("Artists results")

	// Row 1 = selection list
	// Row 2 = selected playlist details
	grid.Clear().AddItem(list, 0, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(grid)

	if len(artists) == 1 {
		SelectArtist(v, artists[0])
	}
}

func CreateSongsPage(v *models.View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Songs")
	v.Pages.AddPage(SONGS_PAGE, box, true, false)
}

func CreateAlbumsPage(v *models.View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Albums")
	v.Pages.AddPage(SONGS_PAGE, box, true, false)
}

func SelectArtist(v *models.View, artist models.SimpleIdentifier) {
	v.UpdateMessageBar(fmt.Sprintf("%s Selected artist %s %s", time.Now().String(), artist.Id, artist.Name))
	pageName, item := v.Pages.GetFrontPage()

	if pageName != ARTISTS_PAGE {
		panic(fmt.Sprintf("This method should only be called from %s, but was called from %s", ARTISTS_PAGE, pageName))
	}

	grid, ok := item.(*tview.Grid)
	if !ok {
		panic(fmt.Sprintf("Expected type *tview.Grid but got %T", grid))
	}
	database, _ := sql.Open("sqlite3", v.DB)
	sqlRows, err := database.Query(
		// TODO: modify the query to also include albums from the tracks that the artist appears on
		"select id, name, total_tracks, release_date, album_type from Album a join AlbumArtist AA on a.id = AA.album_id where AA.artist_id = @Id",
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

	grid.AddItem(table, 1, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(table)
}

func SelectSong(v *models.View, id string, name string) {

}

func SelectAlbum(v *models.View, id string, name string) {

}
