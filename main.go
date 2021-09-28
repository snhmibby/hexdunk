package main

import (
	G "github.com/AllenDang/giu"
)

func draw() {
	G.SingleWindowWithMenuBar().Layout(
		mkMenuWidget(),
		//makeToolBar(),
		mkTabWidget(),
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
		OpenHexFile(n)
	}

	G.SetDefaultFont("DejavuSansMono.ttf", 14)
	w := G.NewMasterWindow("HexDunk", 800, 800, 0)
	w.Run(draw)
}
