package ui

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/storage"
)

var Version string

func Start() {
	client := &http.Client{}
	db, err := storage.ConnectDB()
	if err != nil {
		panic("err when trying to connect the database")
	}
	p := tea.NewProgram(NewModel(client, db))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
