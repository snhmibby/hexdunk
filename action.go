package main

//editor user actions that touch/need the activefile.

import (
	"fmt"
	"io"
	"os"
)

func actionRedo() {
	file := ActiveFile()
	if file == nil {
		return
	}
	file.Redo()
}

func actionUndo() {
	file := ActiveFile()
	if file == nil {
		return
	}
	file.Undo()
}

func actionInsert(b byte) {
	tab := ActiveTab()
	file := ActiveFile()
	if tab == nil || file == nil {
		return
	}
	off := tab.view.cursor
	file.buf.Insert1(off, b)
	tab.view.cursor++
	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() {
			file.Cut(off, 1)
			tab.view.cursor = off
		},
		redo: func() {
			file.buf.Insert1(off, b)
			tab.view.cursor = off
		},
	})
}

func actionOverWrite(b byte) {
	tab := ActiveTab()
	file := ActiveFile()
	if tab == nil || file == nil {
		return
	}
	off := tab.view.cursor
	_, err := file.buf.Seek(off, io.SeekStart)
	if err != nil {
		return //TODO this is not error 'handling'
	}
	var overwritten_ = make([]byte, 1)
	_, err = file.buf.Read(overwritten_)
	if err != nil {
		return //TODO: as above
	}
	overwritten := overwritten_[0]
	file.buf.Remove(off, 1)
	file.buf.Insert1(off, b)
	tab.view.cursor++

	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() {
			file.buf.Remove(off, 1)
			file.buf.Insert1(off, overwritten)
			tab.view.cursor = off
		},
		redo: func() {
			file.buf.Remove(off, 1)
			file.buf.Insert1(off, b)
			tab.view.cursor = off + 1
		},
	})
}

func actionCut() {
	tab := ActiveTab()
	file := ActiveFile()
	if tab == nil || file == nil {
		return
	}
	off, size := tab.view.Selection()
	cut, err := file.Cut(off, size)
	if err != nil {
		ErrorDialog(fmt.Sprintf("Error in action: Cut(%d, %d)", off, size), fmt.Sprint(err))
		return
	}

	HD.ClipBoard = cut
	tab.view.cursor = off
	tab.view.SetSelection(off, 0)
	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() {
			file.buf.Paste(off, cut)
			tab.view.cursor = off
		},
		redo: func() {
			file.buf.Cut(off, size)
			tab.view.cursor = off
		},
	})
}

func actionCopy() {
	file := ActiveFile()
	tab := ActiveTab()
	if tab == nil || file == nil {
		return
	}
	off, size := tab.view.Selection()
	cpy, err := file.Copy(off, size)
	if err != nil {
		ErrorDialog(fmt.Sprintf("Error in action: Copy(%d, %d)", off, size), fmt.Sprint(err))
	}
	HD.ClipBoard = cpy
	tab.view.cursor = off
	tab.view.SetSelection(off, 0)
	file.emptyRedo()
}

//paste in front cursor
func actionPaste() {
	file := ActiveFile()
	tab := ActiveTab()
	if HD.ClipBoard == nil || file == nil || tab == nil {
		return
	}
	off := tab.view.cursor
	buf := HD.ClipBoard //XXX this creates hidden copies of a file-based tree
	file.Paste(off, buf)

	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() {
			file.buf.Cut(off, buf.Size())
			tab.view.cursor = off
		},
		redo: func() {
			file.buf.Paste(off, buf)
			tab.view.cursor = off
		},
	})
}

func actionNewFile() {
	tmpPath, err := os.CreateTemp("", "NewFile*")
	if err != nil {
		ErrorDialog("NewFile", "Cannot create tmp file")
	}
	_, err = OpenHexFile(tmpPath.Name())
	if err != nil {
		ErrorDialog("NewFile", fmt.Sprintf("Couldn't open new tempfile %s", tmpPath.Name()))
	}
}

//callbacks for dialogs are set in the draw() layout function
func dialogOpenCB(p string) {
	_, err := OpenHexFile(p)
	if err != nil {
		title := fmt.Sprintf("Opening File <%s>", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}
}

func actionOpenFile() {
	OpenFileDialog(DialogOpen)
}

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

	//TODO XXX open a new buffer on the whole file again here?
	//this would refresh the working tree buffer
}

func actionSaveFile() {
	OpenFileDialog(DialogSaveAs)
}

func actionSaveAs() {
	OpenFileDialog(DialogSaveAs)
}

func actionCloseTab() {
	if HD.ActiveTab >= 0 {
		CloseTab(HD.ActiveTab)
	}
}

func actionQuit() {
	//TODO: "do you want to save unsaved changes..." dialog
	os.Exit(0)
}
