package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/ui/components"
	"github.com/jeisaraja/youmui/ui/keys"
)

type Tabs struct {
	tabList       []TabItem
	selectTab     int
	content_focus bool
}

type TabItem struct {
	name string
	item components.ContentModel
}

func NewTabs(tabList []TabItem) *Tabs {
	return &Tabs{
		selectTab: 0,
		tabList:   tabList,
	}
}

func (t *Tabs) Init() tea.Cmd { return nil }

func (t *Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			t.decrementSelection()
		case key.Matches(msg, keys.Down):
			t.incrementSelection()
		}
	}
	return t, nil
}

func (t *Tabs) View() string {
	renderedTabs := make([]string, 0)
	renderedTabs = append(renderedTabs, header())

	for i, tl := range t.tabList {
		if i == t.selectTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tl.name))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(tl.name))
		}
	}

	return tabGroupStyle.Render(lipgloss.JoinVertical(lipgloss.Right, renderedTabs...))
}

func (t *Tabs) CurrentTab() TabItem {
	return t.tabList[t.selectTab]
}

func (t *Tabs) incrementSelection() {
	if t.selectTab == len(t.tabList)-1 {
		t.selectTab = 0
	} else {
		t.selectTab++
	}
}

func (t *Tabs) decrementSelection() {
	if t.selectTab > 0 {
		t.selectTab--
	} else {
		t.selectTab = len(t.tabList) - 1
	}
}

func (t *Tabs) SetTab(name string) {
	for i, tab := range t.tabList {
		if tab.name == name {
			t.selectTab = i
		}
	}
}
