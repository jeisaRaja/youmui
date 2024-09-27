package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var Version string

func Start(){
  p:= tea.NewProgram(NewModel())
  if err := p.Start(); err != nil {
    fmt.Println("Error running program: ", err)
    os.Exit(1)
  }
}
