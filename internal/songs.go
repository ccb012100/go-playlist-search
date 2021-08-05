package internal

import (
	"fmt"

	"github.com/ccb012100/go-playlist-search/internal/data"
	"github.com/ccb012100/go-playlist-search/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SearchForSongs(v *models.View) {
	box := tview.NewBox().SetBorder(true).SetTitle("Songs")
	box.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyESC:
			GoToMainMenu(v)
		}

		return e
	})
	v.SetMainPanel(box)
}

func SelectSong(v *models.View, id string, name string) {

}

func ShowDuplicateSongsinStarredPlaylists(v *models.View) {
	v.UpdateMessageBar("func ShowDuplicateSongsinStarredPlaylists()")

	duplicates := data.GetDuplicateTracksInStarredPlaylists(v.DB)

	v.UpdateTitleBar(fmt.Sprintf("%d duplicate songs in Starred Playlists", len(duplicates)))

	displayDuplicateSongs(v, duplicates)
}

func displayDuplicateSongs(v *models.View, dupes []models.DuplicateTrack) {
	table := tview.NewTable().SetBorders(true)
	// set header row
	table.SetCell(0, 0, tview.NewTableCell("Track Name").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 1, tview.NewTableCell("Artists").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 2, tview.NewTableCell("Album Name").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))
	table.SetCell(0, 3, tview.NewTableCell("Playlists").SetTextColor(tcell.ColorOrange).SetAlign(tview.AlignCenter).SetExpansion(1))

	// set table contents
	// start row at 1 to offset for header
	for i := 0; i < len(dupes); i++ {
		dupe := dupes[i]

		// use i+1 to offset for header row
		table.SetCell(i+1, 0, tview.NewTableCell(padLeft(dupe.TrackName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 1, tview.NewTableCell(padLeft(dupe.Artists)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 2, tview.NewTableCell(padLeft(dupe.AlbumName)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(i+1, 3, tview.NewTableCell(padLeft(dupe.Playlists)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	table.SetInputCapture(BackToViewListFunc(v))

	v.SetMainPanel(table)
}
