package data

import (
	"database/sql"
	"sort"

	"github.com/ccb012100/go-playlist-search/internal/models"
)

func GetAlbumsByArtist(artist *models.SimpleIdentifier, db string) []models.Album {
	database, _ := sql.Open("sqlite3", db)

	// Get albums by the artist
	albumArtistRows, err := database.Query(
		"select id, name, total_tracks, release_date, album_type from Album a join AlbumArtist AA on a.id = AA.album_id where AA.artist_id = @Id",
		sql.Named("Id", artist.Id))

	if err != nil {
		panic(err)
	}

	var albums []models.Album
	// track albums in a map so that we display a unique set
	var set = make(map[string]bool)

	for albumArtistRows.Next() {
		var id, name, releaseDate, albumType string
		var totalTracks int

		if err := albumArtistRows.Scan(&id, &name, &totalTracks, &releaseDate, &albumType); err != nil {
			panic(err)
		}

		// skip if the album is already in the slice
		if _, ok := set[id]; ok {
			continue
		}

		set[id] = true

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

		// skip if the album is already in the slice
		if _, ok := set[id]; ok {
			continue
		}

		set[id] = true

		albums = append(albums, models.Album{
			Id:          id,
			Name:        name,
			TotalTracks: totalTracks,
			ReleaseDate: releaseDate,
			AlbumType:   albumType,
		})
	}

	// sort albums
	sort.Sort(models.ByReleaseDate(albums))

	return albums
}

func SearchArtists(query string, db string) []models.SimpleIdentifier {
	database, _ := sql.Open("sqlite3", db)

	rows, err := database.Query("SELECT id, name FROM Artist WHERE name LIKE '%' || @Query || '%' ORDER BY name",
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

	return artists
}

func FindPlaylistsContainingArtist(artist models.SimpleIdentifier, db string) []models.SimpleIdentifier {
	database, _ := sql.Open("sqlite3", db)
	sqlRows, err := database.Query(
		"select PL.id, PL.name from Playlist PL join PlaylistTrack PT on PL.id = PT.playlist_id join Track T on PT.track_id = T.id join TrackArtist TA on T.id = TA.track_id where TA.artist_id = @Id group by PL.id, PL.name order by Pl.name",
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

	return playlists
}

func SearchStarredPlaylists(query string, db string) []models.StarredPlaylistMatch {
	database, _ := sql.Open("sqlite3", db)

	// Get albums by the artist
	sqlRows, err := database.Query(
		"SELECT P.name AS playlistName, T.name AS trackName, A.name AS albumName, GROUP_CONCAT(A2.name, '; ') AS artists FROM Playlist P JOIN PlaylistTrack PT ON P.id = PT.playlist_id JOIN Track T ON PT.track_id = T.id JOIN Album A ON T.album_id = A.id JOIN TrackArtist TA ON T.id = TA.track_id JOIN Artist A2 ON TA.artist_id = A2.id WHERE P.name LIKE 'Starred%' AND (A2.name LIKE '%' || @Query || '%' OR T.name LIKE '%' || @Query || '%' OR A.name LIKE '%' || @Query || '%') GROUP BY P.name, T.id, A.id, PT.added_at, T.track_number ORDER BY P.name, A.id, PT.added_at, T.track_number",
		sql.Named("Query", query))

	if err != nil {
		panic(err)
	}

	var matches []models.StarredPlaylistMatch

	for sqlRows.Next() {
		var playlistName, trackName, albumName, artists string

		if err := sqlRows.Scan(&playlistName, &trackName, &albumName, &artists); err != nil {
			panic(err)
		}

		matches = append(matches, models.StarredPlaylistMatch{
			PlaylistName: playlistName,
			TrackName:    trackName,
			AlbumName:    albumName,
			Artists:      artists,
		})
	}

	return matches
}

func SearchPlaylists(query string, db string) []models.SimpleIdentifier {
	database, _ := sql.Open("sqlite3", db)
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

	return playlists
}
