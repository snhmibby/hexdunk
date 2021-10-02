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

func menuFileOpen() {
	OpenFileDialog(DialogOpen)
}

func menuFileSave() {
	OpenFileDialog(DialogSaveAs)
}
func menuFileSaveAs() {
	OpenFileDialog(DialogSaveAs)
}

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

func menuFile() G.Widget {
	return G.Menu("File").Layout(
		G.MenuItem("New").OnClick(menuFileNew),
		G.MenuItem("Open").OnClick(menuFileOpen),
		G.Separator(),
		G.MenuItem("Save").OnClick(menuFileSave),
		G.MenuItem("Save As").OnClick(menuFileSaveAs),
		G.MenuItem("Close Tab").OnClick(menuFileCloseTab),
		G.Separator(),
		G.MenuItem("Settings").OnClick(menuEditPreferences),
		G.Separator(),
		G.MenuItem("Quit").OnClick(menuFileQuit),
	)
}

func menuEdit() G.Widget {
	return G.Menu("Edit").Layout(
		G.MenuItem("Cut").OnClick(menuEditCut),
		G.MenuItem("Copy").OnClick(menuEditCopy),
		G.MenuItem("Paste").OnClick(menuEditPaste),
	)
}

func mkMenu() G.Widget {
	return G.MenuBar().Layout(
		menuFile(),
		menuEdit(),
		G.Menu("Plugin").Layout(
			G.MenuItem("Load"),
			G.Separator(),
			G.Menu("Plugin Foo").Layout(
				G.MenuItem("DoBar"),
				G.MenuItem("DoBaz"),
				G.MenuItem("Quux"),
				G.Separator(),
				G.MenuItem("Settings"),
			),
		),
	)
}
