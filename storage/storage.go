package storage

import (
	"database/sql"
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
            id TEXT PRIMARY KEY,
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

func AddPlaylist(db *sql.DB, title string) int64 {
	result, err := db.Exec("INSERT INTO playlists (title) VALUES (?)", title)
	if err != nil {
		log.Fatal(err)
	}
	playlistID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return playlistID
}

func AddSongToPlaylist(db *sql.DB, playlistID int64, song api.Song) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT OR IGNORE INTO songs (id, title, description, channel_id, url)
		VALUES (?, ?, ?, ?, ?)`, song.ID, song.Title, song.Description, song.ChannelID, song.URL)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO playlist_songs (playlist_id, song_id)
		VALUES (?, ?)`, playlistID, song.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
