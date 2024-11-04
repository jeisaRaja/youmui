package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeisaraja/youmui/storage"
	"github.com/jeisaraja/youmui/ui"
)

func main() {
	dbPath, err := getDBPath()
	if err != nil {
		panic("error when trying to get the db path")
	}
	db, err := storage.ConnectDB(*dbPath)
	if err != nil {
		panic("[ERROR] when trying to connect the database")
	}
	ui.Start(db)
}

func getDBPath() (*string, error) {
	var dbPath string
	var customPath = "/.local/share/youmui/youmui.db"
	dirName, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user home directory: %w", err)
	}

	dbPath = filepath.Join(dirName, customPath)

	return &dbPath, nil
}
