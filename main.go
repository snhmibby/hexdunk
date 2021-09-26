package hexdunk

import (
	G "github.com/AllenDang/giu"
	//I "github.com/AllenDang/imgui-go"
	B "github.com/snhmibby/filebuf"
)

//for testing

var fileName string
var openBuffer *B.FileBuffer

func draw() {
	hv := HexView("hexview##1", openBuffer)
	G.SingleWindow().Layout(
		G.Label(fileName),
		G.Child().Layout(hv),
	)
}

func main() {
	var err error
	fileName = "go.mod"
	openBuffer, err = B.NewFileBuffer(fileName)
	if err != nil {
		panic(err)
	}
	G.SetDefaultFont("DejavuSansMono.ttf", 11)
	w := G.NewMasterWindow("HexViewer", 800, 600, 0)
	w.Run(draw)
}
