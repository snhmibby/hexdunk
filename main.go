package main

import (
	H "hexdunk/widget"
	"log"

	G "github.com/AllenDang/giu"

	B "github.com/snhmibby/filebuf"
	//I "github.com/AllenDang/imgui-go"
)

type openFile struct {
	name string
	buf  *B.FileBuffer
}

var Tabs []openFile

func makeMenuBar() *G.MenuBarWidget {
	return G.MenuBar().Layout(
		G.Menu("File").Layout(
			G.MenuItem("New"),
			G.MenuItem("Open"),
			G.MenuItem("Save"),
		),
		G.Menu("Edit").Layout(
			G.MenuItem("Cut"),
			G.MenuItem("Copy"),
			G.MenuItem("Paste"),
			G.Separator(),
			G.MenuItem("Preferences"),
		),
	)
}

func draw() {
	var tabs [](*G.TabItemWidget)
	for _, t := range Tabs {
		//TODO: create a nice name... i.e. filename(/a/really/long/directory/path/file.foo) -> file.foo
		var item *G.TabItemWidget = G.TabItem(t.name).Layout(H.HexView("hexview##"+t.name, t.buf))
		tabs = append(tabs, item)
	}
	G.SingleWindowWithMenuBar().Layout(
		makeMenuBar(),
		//makeToolBar(),
		G.TabBar().TabItems(tabs...),
	)
}

//testfiles for now:
var testfiles = []string{
	"go.mod",
	"main.go",
	"/home/jurjen/movies/films/The Last Trapper.mp4",
	"/home/jurjen/Downloads/sas-red-notice-1080p.MP4",
	"/home/jurjen/movies/films/Puss in Boots[2011]BRRip XviD-ETRG/Puss in Boots[2011]BRRip XviD-ETRG.avi",
	"TODO",
	"./README",
}

func main() {
	for _, n := range testfiles {
		buf, err := B.NewFileBuffer(n)
		if err != nil {
			log.Fatal(n, err)
		}
		Tabs = append(Tabs, openFile{n, buf})
	}

	G.SetDefaultFont("DejavuSansMono.ttf", 12)
	w := G.NewMasterWindow("HexViewer", 800, 600, 0)
	w.Run(draw)
}
