package main

import "testing"

func TestGetUserDir(t *testing.T) {
	expectedPath := "/home/jeisa/.local/share/youmui/youmui.db"

	dbPath, err := getDBPath()
	if err != nil {
		t.Fatalf("failed to get the db path: %v", err)
	}

	if *dbPath != expectedPath {
		t.Fatalf("getDBPath is not as expected, got %s, want %s", *dbPath, expectedPath)
	}
}
