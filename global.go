package main

//global variables for hexdunk

import (
	"fmt"
	hexdunk "hexdunk/widget"
	"io/fs"

	B "github.com/snhmibby/filebuf"
	//I "github.com/AllenDang/imgui-go"
)

//global constants
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
	view *hexdunk.HexViewState
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

func mkErr(msg string, e error) error {
	return fmt.Errorf("%s: %s: %v", ProgramName, msg, e)
}
