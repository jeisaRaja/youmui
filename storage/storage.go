package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jeisaraja/youmui/api"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "youmui.db")
	if err != nil {
		return nil, err
	}

	createTables(db)

	return db, nil
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS playlists (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS songs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            channel_id TEXT,
            url TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS playlist_songs (
            playlist_id INTEGER,
            song_id TEXT,
            FOREIGN KEY (playlist_id) REFERENCES playlists(id),
            FOREIGN KEY (song_id) REFERENCES songs(id),
            PRIMARY KEY (playlist_id, song_id)
        );
    `)
	if err != nil {
		log.Fatal(err)
	}
}

func AddPlaylist(db *sql.DB, title string) (*int64, error) {
	result, err := db.Exec("INSERT INTO playlists (title) VALUES (?)", title)
	if err != nil {
		return nil, err
	}
	playlistID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &playlistID, nil
}

func GetPlaylists(db *sql.DB) ([]struct {
	Title string
	ID    int64
	Count int64
}, error) {
	rows, err := db.Query(`
        SELECT 
            playlists.id, 
            playlists.title, 
            COUNT(playlist_songs.song_id) AS song_count
        FROM 
            playlists
        LEFT JOIN 
            playlist_songs ON playlists.id = playlist_songs.playlist_id
        GROUP BY 
            playlists.id, playlists.title;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playlists []struct {
		Title string
		ID    int64
		Count int64
	}
	for rows.Next() {
		var playlist struct {
			Title string
			ID    int64
			Count int64
		}
		if err := rows.Scan(&playlist.ID, &playlist.Title, &playlist.Count); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func AddSongToPlaylist(db *sql.DB, playlistID int64, song api.Song) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var songID int64

	err = tx.QueryRow(`
    SELECT id FROM songs WHERE url = ?
    `, song.URL).Scan(&songID)
	if err == sql.ErrNoRows {
		res, err := tx.Exec(`
		INSERT INTO songs ( title,   url)
		VALUES (?, ?)`, song.Title, song.URL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert song: %w", err)
		}
		songID, err = res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to retreive last song id: %w", err)
		}
	} else if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check if song exist: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO playlist_songs (playlist_id, song_id)
		VALUES (?, ?)`, playlistID, songID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert playlist_song: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
