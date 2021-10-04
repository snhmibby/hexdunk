package main

import (
	"fmt"
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
	if HD.ActiveTab >= 0 {
		CloseTab(HD.ActiveTab)
	}
}

func menuFileQuit() {
	//TODO: "do you want to save unsaved changes..." dialog
	os.Exit(0)
}

//XXX TODO: these should do more error (size/offset, among others) checks!
func menuEditCut() {
	file := ActiveFile()
	tab := ActiveTab()
	st := tab.view
	if file != nil {
		offs, size := st.Selection()
		if size > 0 {
			if st.inSelection(st.cursor) {
				st.cursor = offs
			}
			if offs < 0 || offs+size > file.buf.Size() {
				ErrorDialog("Bad Cut?", fmt.Sprintf("Cut parameters (off: %d, size: %d) are outside of file.\nThis is a bug. How did you get outside of the file? :(", offs, size))
				return
			}
			HD.ClipBoard = file.buf.Cut(offs, size)
			st.SetSelection(offs, 0)
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

func menuEditSettings() {
	//TODO
}

func menuFile() G.Widget {
	return G.Layout{
		G.MenuItem("New").OnClick(menuFileNew),
		G.MenuItem("Open").OnClick(menuFileOpen),
		G.Separator(),
		G.MenuItem("Save").OnClick(menuFileSave),
		G.MenuItem("Save As").OnClick(menuFileSaveAs),
		G.MenuItem("Close Tab").OnClick(menuFileCloseTab),
		G.Separator(),
		//G.MenuItem("Settings").OnClick(menuEditSettings),
		//G.Separator(),
		G.MenuItem("Quit").OnClick(menuFileQuit),
	}
}

func menuEdit() G.Widget {
	return G.Layout{
		G.MenuItem("Cut").OnClick(menuEditCut),
		G.MenuItem("Copy").OnClick(menuEditCopy),
		G.MenuItem("Paste").OnClick(menuEditPaste),
	}
}

func menuPlugin() G.Widget {
	return G.Layout{
		G.MenuItem("Load"),
		G.Separator(),
		G.Menu("Plugin Foo").Layout(
			G.MenuItem("DoBar"),
			G.MenuItem("DoBaz"),
			G.MenuItem("Quux"),
			G.Separator(),
			G.MenuItem("Settings"),
		),
	}
}

func mkMenu() G.Widget {
	return G.MenuBar().Layout(
		G.Menu("File").Layout(menuFile()),
		G.Menu("Edit").Layout(menuEdit()),
		//G.Menu("Plugin").Layout(menuPlugin()),
	)
}
