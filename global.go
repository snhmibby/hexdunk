package main

//global variables and definitions

import (
	"fmt"
	"io/fs"
	"math"

	I "github.com/AllenDang/imgui-go"
	B "github.com/snhmibby/filebuf"
)

const (
	ProgramName = "HexDunk"

	//dialog ids
	DialogOpen   = "Open"
	DialogSaveAs = "Save As"
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

/* some utility functions */

//little hack for vec2
func vec2Abs(v I.Vec2) float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

//mkErr will create a properly formatted error message
func mkErr(msg string, e error) error {
	err := fmt.Errorf("%s: %v", msg, e)
	return err
}
