package main

import (
	G "github.com/AllenDang/giu"
)

func menuEditSettings() {
	//TODO
}

func menuFile() G.Widget {
	return G.Layout{
		G.MenuItem("New").OnClick(actionNewFile),
		G.MenuItem("Open").OnClick(actionOpenFile),
		G.Separator(),
		G.MenuItem("Save").OnClick(actionSaveFile),
		G.MenuItem("Save As").OnClick(actionSaveAs),
		G.MenuItem("Close Tab").OnClick(actionCloseTab),
		G.Separator(),
		//G.MenuItem("Settings").OnClick(menuEditSettings),
		//G.Separator(),
		G.MenuItem("Quit").OnClick(actionQuit),
	}
}

func menuEdit() G.Widget {
	return G.Layout{
		G.MenuItem("Cut").OnClick(actionCut),
		G.MenuItem("Copy").OnClick(actionCopy),
		G.MenuItem("Paste").OnClick(actionPaste),
		G.Separator(),
		G.MenuItem("Undo").OnClick(actionUndo),
		G.MenuItem("Redo").OnClick(actionRedo),
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
