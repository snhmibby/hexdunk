package main

import (
	"fmt"
	"io"
	"os"

	G "github.com/AllenDang/giu"
)

func dialogOpenCB(p string) {
	_, err := OpenHexFile(p)
	if err != nil {
		title := fmt.Sprintf("Error Opening File <%s>", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}
}

func dialogSaveAsCB(p string) {
	hf := ActiveFile()
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		title := fmt.Sprintf("Error Opening File <%s> (for saving).", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}

	hf.buf.Seek(0, io.SeekStart)
	n, err := io.Copy(f, hf.buf)
	if err != nil || n != hf.buf.Size() {
		title := fmt.Sprintf("Error Writing File <%s>", p)
		msg := fmt.Sprintf("Written %d bytes (expected %d)\nError: %v", n, hf.buf.Size(), err)
		ErrorDialog(title, msg)
	}
}

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
