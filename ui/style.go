package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var theme = ThemeNord()
var tabStyle = lipgloss.NewStyle().
	PaddingLeft(1).
	PaddingRight(1)

var activeTabStyle = tabStyle.
	Bold(true).
	Background(theme.SecondaryColor).
	Foreground(theme.PrimaryColor)

var inactiveTabStyle = tabStyle.
	Bold(false).
	Foreground(theme.SecondaryColor)

var tabGroupStyle = lipgloss.NewStyle().
	MarginRight(1).
	MarginLeft(1).
	PaddingBottom(1).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(theme.PrimaryForeground).
	BorderBottom(true).Width(10)

var queueTabStyle = lipgloss.NewStyle().
	MarginLeft(3).
	BorderLeft(true).
	BorderStyle(lipgloss.NormalBorder()).
	PaddingLeft(3).
	Width(40)

var tabContentStyled = lipgloss.NewStyle().
	Width(40).
	MarginRight(1).
	BorderRight(true).
	PaddingRight(1)
