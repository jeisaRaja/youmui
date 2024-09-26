package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "down", "j", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "up", "k", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}
	return m, nil
}

var (
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabBorder = lipgloss.RoundedBorder()
	activeTabBorder   = lipgloss.RoundedBorder()
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder).BorderForeground(lipgloss.Color("#555555")).Padding(0)
	activeTabStyle    = lipgloss.NewStyle().Border(activeTabBorder).BorderForeground(highlightColor).Padding(0).Bold(true)
	docStyle          = lipgloss.NewStyle().Padding(0)
	windowStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 0).Width(30)
)

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	// Render each tab
	for i, t := range m.Tabs {
		var style lipgloss.Style
		if i == m.activeTab {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	// Join tabs vertically
	tabsColumn := strings.Join(renderedTabs, "\n")

	// Render content next to tabs
	content := windowStyle.Render(m.TabContent[m.activeTab])

	// Join the tabs and content horizontally
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabsColumn, content)

	doc.WriteString(row)
	return docStyle.Render(doc.String())
}

func main() {
	tabs := []string{"Tab1", "Tab2", "Tab3", "Tab4"}
	tabContent := []string{
		"Content for Tab1",
		"Content for Tab2",
		"Content for Tab3",
		"Content for Tab4",
	}
	m := model{Tabs: tabs, TabContent: tabContent}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
