package main

import (
	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
)

func draw() {
	//Imgui cannot query the size?? so just set it and hardcode it for my own scrollbar
	I.PushStyleVarFloat(I.StyleVarScrollbarSize, 20)
	I.PushStyleVarFloat(I.StyleVarScrollbarRounding, 5)
	defer I.PopStyleVarV(2)

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
