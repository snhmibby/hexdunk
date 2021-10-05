package main

import (
	G "github.com/AllenDang/giu"
)

func draw() {
	G.SingleWindowWithMenuBar().Layout(
		G.PrepareMsgbox(),
		PrepareFileDialog(DialogOpen, dialogOpenCB),
		PrepareFileDialog(DialogSaveAs, dialogSaveAsCB),
		mkMenu(),
		//makeToolBar(),
		mkTabWidget(),
	)
}

func main() {
	G.SetDefaultFont("DejavuSansMono.ttf", 12)
	w := G.NewMasterWindow("HexDunk", 800, 800, 0)
	w.Run(draw)
}
