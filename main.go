package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ccb012100/go-playlist-search/config"
	"github.com/ccb012100/go-playlist-search/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/text"
)

const (
	// redrawInterval is how often termdash redraws the screen.
	redrawInterval = 250 * time.Millisecond
	rootId         = "root-ID"
	leftId         = "left-ID"
	rightId        = "right-ID"
)

func main() {
	// application struct
	app := &models.App{}

	conf := config.SetConfig()
	app.DB = conf.DBFilePath

	t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		panic(fmt.Errorf("tcell.New => %v", err))
	}
	// calling t.Close() is necessary to exit with a clean terminal state
	defer t.Close()

	btn, err := button.New("button", func() error {
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("button.New => %v", err))
	}

	txt, err := text.New()
	if err != nil {
		panic(fmt.Errorf("text.New => %v", err))
	}

	msg, err := text.New()
	if err != nil {
		panic(fmt.Errorf("text.Net => %v", err))
	}

	app.MessagePanel = msg

	left := container.Left(
		container.ID(leftId),
		container.SplitHorizontal(
			container.Top(container.PlaceWidget(btn),
				container.PaddingBottom(1),
				container.PaddingTop(1),
				container.PaddingLeft(1),
				container.PaddingRight(1),
				container.Border(linestyle.Round),
				container.BorderColor(cell.ColorGreen),
				container.BorderTitle("Top"),
				container.BorderTitleAlignCenter()),
			container.Bottom(
				container.PlaceWidget(msg),
				container.PaddingBottom(1),
				container.PaddingTop(1),
				container.PaddingLeft(1),
				container.PaddingRight(1),
				container.BorderColor(cell.ColorAqua),
				container.BorderTitle("Messages"),
				container.BorderTitleAlignCenter(),
				container.Border(linestyle.Round)),
			container.SplitPercent(75)),
	)

	right := container.Right(
		container.ID(rightId),
		container.PlaceWidget(txt),
		container.PaddingBottom(1),
		container.PaddingTop(1),
		container.PaddingLeft(1),
		container.PaddingRight(1),
		container.BorderTitle("Right"),
		container.Border(linestyle.Light),
		container.BorderColor(cell.ColorSilver),
		container.BorderTitleAlignRight())

	c, err := container.New(t,
		container.SplitVertical(left, right),
		container.KeyFocusNext(keyboard.KeyCtrlN),
		container.KeyFocusPrevious(keyboard.KeyCtrlP))
	if err != nil {
		panic(fmt.Errorf("container.New => %v", err))
	}

	app.Root = c

	// Context for Termdash, the application quits when this expires.
	// Since the context has no deadline, it will only expire when cancel() is called.
	ctx, cancel := context.WithCancel(context.Background())

	// A keyboard subscriber that terminates the application by cancelling the context.
	// quit app by typing 'q', 'ESC' or 'Ctrl-C'
	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}

	app.MessagePanel.Write("App started")

	// Run app
	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}
}
