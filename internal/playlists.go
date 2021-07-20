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
