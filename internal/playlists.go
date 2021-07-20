package pages

import (
	"database/sql"
	"fmt"

	"github.com/ccb012100/go-playlist-search/internal/models"
	"github.com/rivo/tview"
)

func ShowPlaylists(v *models.View, query string) {
	v.UpdateMessageBar(fmt.Sprintf("func ShowPlaylists() query='%s'", query))
	pageName, _ := v.Pages.GetFrontPage()

	if pageName != PLAYLISTS_PAGE {
		panic(fmt.Sprintf("This method should only be called from %s, but was called from %s", PLAYLISTS_PAGE, pageName))
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

	AddQuitToHomeOption(list, v)

	AddResetOption(list, func() {
		CreatePlaylistsPage(v)
		v.Pages.ShowPage(PLAYLISTS_PAGE)
		v.App.SetFocus(v.Pages)
	})

	list.SetTitle("Playlists")

	// Row 1 = selection list
	// Row 2 = selected playlist details
	v.Grid.AddItem(list, 1, 0, 1, 1, 0, 0, true)
	v.App.SetFocus(list)
}

func SelectPlaylist(v *models.View, playlist models.SimpleIdentifier) {
	pageName, _ := v.Pages.GetFrontPage()

	if pageName != PLAYLISTS_PAGE {
		panic("This method should only be called from the Playlists page")
	}

	var textView = tview.NewTextView()

	textView.SetText(fmt.Sprintf("Selected playlist id='%s', name='%s'", playlist.Id, playlist.Name)).
		SetTitle(fmt.Sprintf("Selected Playlist: %s", playlist.Name))

	v.Grid.AddItem(textView, 1, 0, 1, 1, 0, 0, false)
}
