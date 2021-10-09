package main

//global variables and definitions

import (
	"fmt"
	"io/fs"

	B "github.com/snhmibby/filebuf"
)

const (
	ProgramName = "HexDunk"

	//dialog ids
	DialogOpen   = "Open"
	DialogSaveAs = "Save As"
)

//probably should not be an opaque buffer, but a struct with an action-enum.
//that way, we can adjust copy-buffers when a file gets saved (i.e. referenced
//portion of the saved file should be removed through all buffers throughout
//the program on a save.)
type Undo struct {
	undo, redo func()
}

//an opened file
type HexFile struct {
	name       string
	buf        *B.Buffer
	dirty      bool
	stats      fs.FileInfo
	undo, redo []Undo
}

//each tab is a view on an opened file
type HexTab struct {
	name string
	view *ViewState
}

type ViewState struct {
	cursor   int64 //address (byte offset in file)
	editmode editMode

	//selected bytes
	selectionStart, selectionSize int64

	//selection dragging
	dragging  bool
	dragstart int64
}

type editMode int

const (
	NormalMode editMode = iota
	InsertMode
	OverwriteMode
)

//global variables
type Globals struct {
	// All opened files, index by file-path
	Files map[string]*HexFile

	// All tabs (every tab is a view on an opened file)
	Tabs []HexTab

	//Index of active tab in display
	ActiveTab int

	//Current copy/paste buffer
	//TODO: could be something nice, a circular buffer, named buffers, etc
	ClipBoard *B.Buffer
}

var HD Globals = Globals{
	Tabs:      make([]HexTab, 0),
	ActiveTab: -1,
	Files:     make(map[string]*HexFile),
}

func ActiveTab() *HexTab {
	if HD.ActiveTab >= 0 {
		//consistency check
		if ActiveFile() == nil {
			panic("impossible")
		}
		return &HD.Tabs[HD.ActiveTab]
	}
	return nil
}

func ActiveFile() *HexFile {
	if HD.ActiveTab >= 0 {
		hf, ok := HD.Files[HD.Tabs[HD.ActiveTab].name]
		if !ok {
			panic("tab opened on closed file")
		}
		return hf
	}
	return nil
}

/* ViewState methods */

func (view *ViewState) SetSelection(begin, size int64) {
	view.selectionStart, view.selectionSize = begin, size
}

func (view *ViewState) Selection() (begin, size int64) {
	return view.selectionStart, view.selectionSize
}

func (st *ViewState) inSelection(addr int64) bool {
	off, size := st.Selection()
	return addr >= off && addr < off+size
}

/* some utility functions */

//mkErr will create a properly formatted error message
func mkErr(msg string, e error) error {
	err := fmt.Errorf("%s: %v", msg, e)
	return err
}
