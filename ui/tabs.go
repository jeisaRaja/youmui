package ui


type Tabs struct {
	tabList   []TabItem
	selectTab int
}

type TabItem struct {
	name string
	item string
}

func NewTabs(tabList []TabItem) *Tabs {
	return &Tabs{
		selectTab: 0,
		tabList:   tabList,
	}
}
