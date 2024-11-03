package ui

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var Version string

func Start(db *sql.DB) {
	client := &http.Client{}
	p := tea.NewProgram(NewModel(client, db))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
