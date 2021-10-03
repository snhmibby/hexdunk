package main

//little program to display a file-system tree and basic info
//once a file-path is selected, print it on stdout

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
)

const (
	timeFmt     = "02 Jan 06 15:04"
	nodeFlags   = I.TreeNodeFlagsSpanFullWidth | I.TreeNodeFlagsOpenOnArrow | I.TreeNodeFlagsOpenOnDoubleClick
	leafFlags   = I.TreeNodeFlagsLeaf
	tableFlags  = I.TableFlags_ScrollX | I.TableFlags_ScrollY | I.TableFlags_Resizable | I.TableFlags_SizingStretchProp
	selectFlags = I.SelectableFlagsAllowDoubleClick | I.SelectableFlagsSpanAllColumns
)

type fileDialog struct {
	//imgui/giu and popup related
	id       string //giu/imgiu id
	open     bool
	callback func(path string)

	//functionality related
	statCache       map[string]fs.FileInfo
	dirCache        map[string][]fs.FileInfo
	showHiddenFiles bool
	selectedFile    string //full path of file selected in fileTable
	currentDir      string //full path of directory selected in dirTree
	startDir        string //starting directory arg or cwd()
}

var _ G.Disposable = &fileDialog{}

func (d *fileDialog) Dispose() { /* empty */ }

func mkSize(sz_ int64) string {
	sizes := []string{"KB", "MB", "GB", "TB"}
	sz := float64(sz_)
	add := ""
	for _, n := range sizes {
		if sz < 1024 {
			break
		}
		sz = sz / 1024
		add = n
	}
	if add == "" {
		return fmt.Sprint(sz_)
	} else {
		return fmt.Sprintf("%.2f %s", sz, add)
	}
}

//statFile follows symbolic links
func (fd *fileDialog) statFile(path string) (fs.FileInfo, error) {
	st, ok := fd.statCache[path]
	if !ok {
		var err error
		st, err = os.Stat(path)
		if err != nil {
			return nil, err
		}
		fd.statCache[path] = st
	}
	return st, nil
}

//return a list of directory entries
func (fd *fileDialog) readDir(path string) ([]fs.FileInfo, error) {
	entry, ok := fd.dirCache[path]
	if !ok {
		direntry, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		entry := make([]fs.FileInfo, len(direntry))
		for i, f := range direntry {
			childPath := filepath.Join(path, f.Name())

			if f.Type()&fs.ModeSymlink == 0 {
				entry[i], err = f.Info()
			} else {
				//follow symlink
				entry[i], err = fd.statFile(childPath)
			}
			if err != nil {
				return nil, err
			}
		}
		fd.dirCache[path] = entry
	}
	return entry, nil
}

//Get all necessary info to display a directory, including giu.TreeNodeFlags*
//any errors will result in (silently) not displaying this directory!!
//path -> (TreeNodeFlags, fileinfo, [children fileinfo], succes/display?)
func (fd *fileDialog) getDirInfo(path string) (int, fs.FileInfo, []fs.FileInfo, bool) {
	info, err := fd.statFile(path)
	if err != nil {
		return 0, nil, nil, false
	}

	entries, err := fd.readDir(path)
	if err != nil {
		return 0, nil, nil, false
	}

	if info.Name()[0] == '.' && !fd.showHiddenFiles {
		return 0, nil, nil, false
	}

	flags := leafFlags
	for _, e := range entries {
		if e.IsDir() {
			flags = nodeFlags
			break
		}
	}

	if path == fd.currentDir {
		flags |= I.TreeNodeFlagsSelected
	}

	return flags, info, entries, true
}

func (fd *fileDialog) dirTree(path string) {
	flags, info, entries, ok := fd.getDirInfo(path)
	if !ok {
		return
	}

	if strings.HasPrefix(fd.startDir, path) {
		flags |= I.TreeNodeFlagsDefaultOpen
	}

	I.PushStyleVarFloat(I.StyleVarIndentSpacing, 5)
	defer I.PopStyleVar()

	open := I.TreeNodeV(info.Name(), flags)
	if path == fd.currentDir {
		I.SetItemDefaultFocus()
	}
	if I.IsItemClicked(int(G.MouseButtonLeft)) {
		fd.currentDir = path
	}

	if open {
		defer I.TreePop()
		for _, e := range entries {
			if e.IsDir() {
				name := filepath.Join(path, e.Name())
				fd.dirTree(name)
			}
		}
	}
}

func isHidden(entry fs.FileInfo) bool {
	return entry.Name()[0] == '.'
}

func (fd *fileDialog) fileTable() {
	I.Text(fd.currentDir)
	if I.BeginTable("FSTable", 3, tableFlags, I.ContentRegionAvail(), 0) {
		defer I.EndTable()
		I.TableSetupColumn("Name", 0, 10, 0)
		I.TableSetupColumn("Size", 0, 2, 0)
		I.TableSetupColumn("Time", 0, 4, 0)
		I.TableSetupScrollFreeze(1, 1)
		I.TableHeadersRow()
		//TODO: set up sorting

		entries, err := fd.readDir(fd.currentDir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.IsDir() || isHidden(e) {
				continue
			}
			path := filepath.Join(fd.currentDir, e.Name())

			I.TableNextRow(0, 0)

			I.TableNextColumn()
			if I.SelectableV(e.Name(), path == fd.selectedFile, selectFlags, I.Vec2{}) {
				fd.selectedFile = path
				if I.IsMouseDoubleClicked(int(G.MouseButtonLeft)) {
					fd.selectFile()
				}
			}
			I.TableNextColumn()
			I.Text(mkSize(e.Size()))
			I.TableNextColumn()
			I.Text(e.ModTime().Format(timeFmt))
		}
	}
}

func (fd *fileDialog) selectFile() {
	file := fd.selectedFile
	if !filepath.IsAbs(file) {
		file = filepath.Join(fd.currentDir, file)
	}
	if fd.callback != nil && fd.selectedFile != "" && file != "" {
		fd.callback(file)
	}

	I.CloseCurrentPopup()
	fd.statCache = make(map[string]fs.FileInfo)
	fd.dirCache = make(map[string][]fs.FileInfo)
	fd.startDir, _ = filepath.Abs(".")
	fd.currentDir = fd.startDir
	fd.selectedFile = ""
	fd.saveState()
}

func (fd *fileDialog) saveState() {
	G.Context.SetState(fd.id, fd)
}

func (fd *fileDialog) close() {
	I.CloseCurrentPopup()
}

func (fd *fileDialog) mkNavBar() {
	width, _ := G.GetAvailableRegion()
	G.InputText(&fd.selectedFile).Size(width).Build()
}

//
//Public:
//

func InfoDialog(title, msg string) {
	G.Msgbox("Info:", title+":\n"+msg).Buttons(G.MsgboxButtonsOk)
}

func ErrorDialog(title, msg string) {
	G.Msgbox("ERROR!", "When: "+title+"\n\nError: "+msg).Buttons(G.MsgboxButtonsOk)
}

func OpenFileDialog(id string) {
	fdRaw := G.Context.GetState(id)
	if fdRaw == nil {
		panic("Couldn't open file dialog " + id)
	}
	fd := fdRaw.(*fileDialog)
	fd.open = true
	fd.saveState()
}

func CloseFileDialog(id string) {
	G.CloseCurrentPopup()
}

func PrepareFileDialog(id string, cb func(string)) G.Widget {
	var fd *fileDialog
	dialogRaw := G.Context.GetState(id)
	if dialogRaw == nil {
		start, _ := filepath.Abs(".")
		fd = &fileDialog{
			id:         id,
			statCache:  make(map[string]fs.FileInfo),
			dirCache:   make(map[string][]fs.FileInfo),
			startDir:   start,
			currentDir: start,
			callback:   cb,
		}
		fd.saveState()
	} else {
		fd = dialogRaw.(*fileDialog)
	}

	return G.Custom(func() {
		if fd.open {
			G.OpenPopup(id)
			fd.open = false
		}
		G.SetNextWindowSizeV(600, 300, G.ConditionOnce)
		G.PopupModal(id).Layout(
			G.Custom(fd.mkNavBar),
			G.Custom(func() {
				//use a child frame to block the lists going off-screen
				//at the bottom of the screen is a row of buttons and such
				w, h := G.GetAvailableRegion()
				//adjust for buttonrow height
				_, spacing := G.GetItemSpacing()
				_, padding := G.GetFramePadding()
				_, buttonH := G.CalcTextSize("F")
				h -= buttonH + padding*2 + spacing
				G.Child().Layout(
					G.SplitLayout(G.DirectionHorizontal, true, 200,
						G.Custom(func() { fd.dirTree(filepath.FromSlash("/")) }),
						G.Custom(fd.fileTable),
					),
				).Border(false).Size(w, h).Build()
			}),
			G.Row(
				G.Checkbox("Show Hidden", &fd.showHiddenFiles),
				G.Button("Cancel").OnClick(fd.close),
				G.Button(id).OnClick(fd.selectFile),
			),
		).Flags(G.WindowFlagsNone).Build()

	})
}
