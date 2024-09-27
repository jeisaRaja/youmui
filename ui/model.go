package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tabs *Tabs
}

var (
	tabItems = map[string]string{
		"Home":     "Main Area",
		"Browse":   "Browse Musics",
		"Playlist": "My Playlists",
		"Search":   "Search Results",
		"Library":  "My Library",
	}
)

func NewModel() *model {
	var m model
	var tabList []TabItem
	for name, item := range tabItems {
		tabList = append(tabList, TabItem{name: name, item: item})
	}
	m.tabs = NewTabs(tabList)
	return &m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k": // Up key
			if m.tabs.selectTab > 0 {
				m.tabs.selectTab--
			}
		case "j": // Down key
			if m.tabs.selectTab < len(m.tabs.tabList)-1 {
				m.tabs.selectTab++
			}
		case "enter": // Select tab
			// Handle selection if needed (currently does nothing)
		case "ctrl+c", "q": // Quit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	var output string
	split := "╠═════════════════════════════╦══════════════════════════╣\n"

	// Header
	output += "╔════════════════════════════════════════════════════════╗\n"
	output += "║ YouTube Music TUI                                        ║\n"
	output += split

	// Sidebar for tabs
	output += "║ Sidebar               ║ Main Area                    ║\n"
	output += "║                       ║                              ║\n"

	for i, tab := range m.tabs.tabList {
		// Format the tab with highlighting for the selected tab
		if i == m.tabs.selectTab {
			output += fmt.Sprintf("║ [%s]                ║ %s   ║\n", tab.name, tab.item)
		} else {
			output += fmt.Sprintf("║ [%s]                ║                              ║\n", tab.name)
		}
	}
	// Add footer and status
	output += split
	output += "║ Status Bar: Playing | Volume: 50%                     ║\n"
	output += "╚════════════════════════════════════════════════════════╝\n"

	return output
}
