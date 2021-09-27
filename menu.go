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

func menuEditCut()         {}
func menuEditCopy()        {}
func menuEditPaste()       {}
func menuEditPreferences() {}

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
			G.MenuItem("Preferences").OnClick(menuEditPreferences),
		),
		G.Menu("View").Layout(),
		G.Menu("Plugin").Layout(),
	)
}
