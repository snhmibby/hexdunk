package main

import (
	G "github.com/AllenDang/giu"
)

func draw() {

	G.SingleWindowWithMenuBar().Layout(
		PrepareFileDialog(DialogOpen, dialogOpenCB),
		mkMenu(),
		//G.Custom(func() { G.OpenPopup(DialogOpen) }),
		//makeToolBar(),
		mkTabWidget(),
	)
}

func main() {
	G.SetDefaultFont("DejavuSansMono.ttf", 12)
	w := G.NewMasterWindow("HexDunk", 800, 800, 0)
	w.Run(draw)
}
