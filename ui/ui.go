package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var Version string

func Start(){
  p:= tea.NewProgram(NewModel())
  if _, err := p.Run(); err != nil {
    fmt.Println("Error running program: ", err)
    os.Exit(1)
  }
}
