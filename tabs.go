package main

import (
	"fmt"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
	//I "github.com/AllenDang/imgui-go"
)

func mkTabWidget() G.Widget {
	//use a child with NoMove flag so the widget can do drag events
	return G.Child().Flags(G.WindowFlagsNoMove).Layout(G.Custom(func() {
		if len(HD.Tabs) != 0 && I.BeginTabBar("HexViewerTabs") {
			defer I.EndTabBar()
			for i, tab := range HD.Tabs {

				hf, ok := HD.Files[tab.name]
				if !ok {
					panic("tab opened on closed file")
				}

				flags := G.TabItemFlagsNone
				if hf.dirty {
					flags |= G.TabItemFlagsUnsavedDocument
				}
				if I.BeginTabItem(hf.stats.Name()) {
					defer I.EndTabItem()
					HD.ActiveTab = i
					h := HexView(fmt.Sprint(i, ".hexview##", hf.name), hf.buf, tab.view)
					h.Build()
				}
			}
		}
	}))
}

func ActiveTab() *HexTab {
	if HD.ActiveTab >= 0 {
		//consistency check
		if ActiveFile() == nil {
			panic("impossible")
		}
		return &HD.Tabs[HD.ActiveTab]
	}
	return nil
}

func ActiveFile() *HexFile {
	if HD.ActiveTab >= 0 {
		hf, ok := HD.Files[HD.Tabs[HD.ActiveTab].name]
		if !ok {
			panic("tab opened on closed file")
		}
		return hf
	}
	return nil
}

func OpenTab(hf *HexFile) {
	HD.Tabs = append(HD.Tabs, HexTab{name: hf.name, view: new(HexViewState)})
}

func CloseTab(t int) {
	if t < 0 || len(HD.Tabs) < t {
		panic("closeTab number doesn't exist (shouldn't happen)")
	}
	tab := HD.Tabs[HD.ActiveTab]
	copy(HD.Tabs[t:], HD.Tabs[t+1:])
	HD.Tabs = HD.Tabs[:len(HD.Tabs)-1]
	if HD.ActiveTab == t {
		HD.ActiveTab--
	}

	//if another tab has the same file open, don't close the file permanently
	for _, t := range HD.Tabs {
		if t.name == tab.name {
			return
		}
	}
	CloseHexFile(tab.name) // we closed last tab on this guy
}
