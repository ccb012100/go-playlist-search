package models

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/text"
)

type App struct {
	// root container
	Root *container.Container
	// widget for displaying messages
	MessagePanel *text.Text
	// db file path
	DB string
}

type Album struct {
	Id          string
	Name        string
	TotalTracks int
	ReleaseDate string
	AlbumType   string
}

type SimpleIdentifier struct {
	Name string
	Id   string
}
