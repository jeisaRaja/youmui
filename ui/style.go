package ui

import "github.com/charmbracelet/lipgloss"

var tabStyle = lipgloss.NewStyle().
	PaddingLeft(1).
	PaddingRight(1)

var activeTabStyle = tabStyle.
	Bold(true)
	// Background(theme.SecondaryColor).
	// Foreground(theme.PrimaryColor)

var inactiveTabStyle = tabStyle.
	Bold(false)
	// Foreground(theme.SecondaryColor)

var tabGroupStyle = lipgloss.NewStyle().
	MarginRight(1).
	MarginLeft(1).
	PaddingBottom(1).
	BorderStyle(lipgloss.NormalBorder()).
	// BorderForeground(theme.PrimaryForeground).
	BorderBottom(true)
