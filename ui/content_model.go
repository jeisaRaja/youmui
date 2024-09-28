package ui

import tea "github.com/charmbracelet/bubbletea"

type content struct {
	msg string
	id  int
	err error
}

type ContentModel interface {
	tea.Model
}
