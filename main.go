package main

import (
	"fmt"
	"io"
	"os"

	G "github.com/AllenDang/giu"
	B "github.com/snhmibby/filebuf"
)

func OpenHexFile(path string) (*HexFile, error) {
	hf, ok := HD.Files[path]
	if !ok {
		//open & stat file
		stats, err := os.Stat(path)
		if err != nil {
			return nil, mkErr("OpenHexFile", err)
		}
		//XXX this check shouldn't even be here - we want to be able to edit ANY file
		if !stats.Mode().IsRegular() {
			return nil, mkErr("OpenHexFile", fmt.Errorf("%s is not a regular file", path))
		}
		buf, err := B.OpenFile(path)
		if err != nil {
			return nil, mkErr("OpenHexFile", err)
		}
		hf = new(HexFile)
		hf.buf = buf
		hf.name = path
		hf.stats = stats
		HD.Files[path] = hf
	}
	OpenTab(hf)
	return hf, nil
}

//should only called when the last view (tab) on this file is closed
func CloseHexFile(path string) error {
	hf, ok := HD.Files[path]
	if !ok {
		return mkErr("CloseHexFile", fmt.Errorf("No file named (%s) open.", path))
	}
	if hf.dirty {
		//TODO: dialog.ReallyClose ? Option to save
	}
	delete(HD.Files, path)

	//sanity check
	for _, t := range HD.Tabs {
		if t.name == path {
			panic("shouldn't happen")
		}
	}
	return nil
}

func dialogOpenCB(p string) {
	_, err := OpenHexFile(p)
	if err != nil {
		title := fmt.Sprintf("Opening File <%s>", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}
}

//XXX This confusingly can change the buffers that were Cut from this file and still need to Read from it
//There are 2 obvious ways to resolve this:
//1) On every Cut/Paste-operation, load all file-blocks into memory
//2) On every save operation, load all affected file-blocks into memory.
//   This would mean that we have to keep track of where they went (were pasted)
func dialogSaveAsCB(p string) {
	hf := ActiveFile()
	f, err := os.CreateTemp("", "")
	if err != nil {
		title := fmt.Sprintf("Opening File <%s> for saving.", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}
	hf.buf.Seek(0, io.SeekStart)
	n, err := io.Copy(f, hf.buf)
	if err != nil || n != hf.buf.Size() {
		os.Remove(f.Name())
		title := fmt.Sprintf("Writing File <%s>", p)
		msg := fmt.Sprintf("Written %d bytes (expected %d)\nError: %v", n, hf.buf.Size(), err)
		ErrorDialog(title, msg)
	}
	err = os.Rename(f.Name(), p)
	if err != nil {
		title := fmt.Sprintf("Naming File <%s>", p)
		msg := fmt.Sprintf("Couldn't rename tmp file <%s> to <%s>", f.Name(), p)
		ErrorDialog(title, msg)
	}

	//XXX open a new buffer on the whole file again here?
	//this would refresh the working tree buffer
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
