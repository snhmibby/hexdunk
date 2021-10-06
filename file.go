package main

import (
	"fmt"
	"os"

	B "github.com/snhmibby/filebuf"
)

func (hf *HexFile) emptyRedo() {
	hf.redo = []Undo{}
}

func (hf *HexFile) addRedo(f Undo) {
	hf.redo = append(hf.redo, f)
}

func (hf *HexFile) addUndo(f Undo) {
	hf.undo = append(hf.undo, f)
}

func (hf *HexFile) Redo() {
	sz := len(hf.redo)
	if sz == 0 {
		return
	}
	f := hf.redo[sz-1]
	hf.redo = hf.redo[:sz-1]
	f.redo()
	hf.addUndo(f)

}

func (hf *HexFile) Undo() {
	sz := len(hf.undo)
	if sz == 0 {
		return
	}
	f := hf.undo[sz-1]
	hf.undo = hf.undo[:sz-1]
	f.undo()
	hf.addRedo(f)
}

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

func (hf *HexFile) Copy(off, size int64) (*B.Buffer, error) {
	if size <= 0 {
		return nil, fmt.Errorf("Cut: size <= 0")
	}
	if off < 0 || off+size > hf.buf.Size() {
		e := fmt.Errorf("Copy: 0 < off (%d) < off + size (%d) < file.Size() (%d)", off, size, hf.buf.Size())
		return nil, e
	}
	return hf.buf.Copy(off, size), nil
}

func (hf *HexFile) Paste(off int64, b *B.Buffer) {
	hf.buf.Paste(off, b)
}

func (hf *HexFile) Cut(off, size int64) (*B.Buffer, error) {
	if size <= 0 {
		return nil, fmt.Errorf("Cut: size <= 0")
	}
	if off < 0 || off+size > hf.buf.Size() {
		e := fmt.Errorf("Cut: 0 < off (%d) < off + size (%d) < file.Size() (%d)", off, size, hf.buf.Size())
		return nil, e
	}
	return hf.buf.Cut(off, size), nil
}
