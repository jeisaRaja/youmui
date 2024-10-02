package types

import tea "github.com/charmbracelet/bubbletea"

type Mockres []string

type FocusLevel int

const (
	TabsFocus = iota
	ContentFocus
	SubContentFocus
)

type FocusMsg struct {
	Level FocusLevel
	Msg   tea.Msg
}
