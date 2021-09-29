package main

import (
	"os"

	G "github.com/AllenDang/giu"
	//I "github.com/AllenDang/imgui-go"
)

func menuFileNew() {
	tmpPath, err := os.CreateTemp("", "NewFile*")
	if err != nil {
		mkErr("menu.File.New: cannot create tmp file", err)
		return
	}
	OpenHexFile(tmpPath.Name())
}

//TODO: dialogs...
func menuFileOpen()   {}
func menuFileSave()   {}
func menuFileSaveAs() {}

func menuFileCloseTab() {
	if HD.ActiveTab < 0 {
		return
	}
	CloseTab(HD.ActiveTab)
}

func menuFileQuit() {
	//TODO: "do you want to save unsaved changes..." dialog
	os.Exit(0)
}

//XXX TODO: these should do more error (size/offset, among others) checks!
func menuEditCut() {
	file := ActiveFile()
	tab := ActiveTab()
	if file != nil {
		offs, size := tab.view.Selection()
		if size > 0 {
			HD.ClipBoard = file.buf.Cut(offs, size)
			tab.view.SetSelection(offs, 0)
		}
	}
}

func menuEditCopy() {
	hf := ActiveFile()
	tab := ActiveTab()
	if hf != nil {
		offs, size := tab.view.Selection()
		if size > 0 {
			HD.ClipBoard = hf.buf.Copy(offs, size)
		}
	}
}

//paste in front cursor
func menuEditPaste() {
	hf := ActiveFile()
	tab := ActiveTab()
	if HD.ClipBoard != nil && hf != nil {
		offs := tab.view.Cursor()
		hf.buf.Paste(offs, HD.ClipBoard)
	}
}

func menuEditPreferences() {
	//TODO
}

var testBool bool

func mkMenuWidget() *G.MenuBarWidget {
	return G.MenuBar().Layout(
		G.Menu("File").Layout(
			G.MenuItem("New").OnClick(menuFileNew),
			G.MenuItem("Open").OnClick(menuFileOpen),
			G.Separator(),
			G.MenuItem("Save").OnClick(menuFileSave),
			G.MenuItem("Save As").OnClick(menuFileSaveAs),
			G.MenuItem("Close Tab").OnClick(menuFileCloseTab),
			G.Separator(),
			G.MenuItem("Quit").OnClick(menuFileQuit),
		),
		G.Menu("Edit").Layout(
			G.MenuItem("Cut").OnClick(menuEditCut),
			G.MenuItem("Copy").OnClick(menuEditCopy),
			G.MenuItem("Paste").OnClick(menuEditPaste),
			G.Separator(),
			G.MenuItem("Settings").OnClick(menuEditPreferences),
		),
		G.Menu("View").Layout(
			G.MenuItem("Search/Replace"),
		),
		G.Menu("Plugin").Layout(
			G.MenuItem("Load"),
			G.Separator(),
			G.Menu("Plugin Foo").Layout(
				G.MenuItem("DoBar"),
				G.MenuItem("DoBaz"),
				G.MenuItem("Quux"),
				G.Separator(),
				G.Checkbox("Activate/Deactivate)", &testBool),
				G.MenuItem("Settings"),
			),
		),
	)
}
