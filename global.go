package main

//global variables and definitions

import (
	"fmt"
	"io/fs"

	B "github.com/snhmibby/filebuf"
)

const (
	ProgramName = "HexDunk"
)

//an opened file
type HexFile struct {
	name  string
	buf   *B.Buffer
	dirty bool
	stats fs.FileInfo
}

//each tab is a view on an opened file
type HexTab struct {
	name string
	view *HexViewState
}

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

/* some utility functions */

//mkErr will create a properly formatted error message
//then it panics with that error message and it is possible to recover from and handle this later?
func mkErr(msg string, e error) error {
	err := fmt.Errorf("%s: %s: %v", ProgramName, msg, e)
	return err
}
