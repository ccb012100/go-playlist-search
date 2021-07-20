package pages

import (
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
