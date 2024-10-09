package ui

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var Version string

func Start() {
	client := &http.Client{}
	p := tea.NewProgram(NewModel(client))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
