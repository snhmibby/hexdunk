package hexdunk

//this kind of ripped is ripped from ocornuts' hexwidget (imgui_hex.cpp)

//TODO: most hex dumps have a printable character view after the raw data, implement it
//TODO: use timestamps to check to see if file is edited by another program

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"unicode"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
	B "github.com/snhmibby/filebuf"
)

func numHexDigits(addr int64) int {
	hexdigits := 0
	for addr > 0 {
		hexdigits++
		addr /= 16
	}
	return hexdigits
}

func addrLabel(addr int64, nDigits int) string {
	return fmt.Sprintf("%0*X: ", nDigits, addr)
}

type TextColor struct {
	fg, bg color.RGBA
}

type cursorAddr struct {
	//byte address of cursor
	addr int64
}

//Widget state
type HexViewState struct {
	//offset in the file where is the cursor (should be on screen?)
	cursor int64

	//selected bytes (offset from cursor)
	selection int64

	//selection dragging
	dragging  bool
	lastmouse int64
}

func (view *HexViewState) Cursor() int64 {
	return view.cursor
}

func (view *HexViewState) SetSelection(addr, size int64) {
	view.cursor, view.selection = addr, size
}

func (view *HexViewState) Selection() (addr, size int64) {
	return view.cursor, view.selection
}

type HexViewWidget struct {
	state *HexViewState

	id     string
	buffer *B.Buffer

	topAddr         int64
	bytesPerLine    int64
	linesPerScreen  int64
	width           float32
	height          float32
	charWidth       float32
	charHeight      float32
	addressBarWidth float32
}

func HexView(id string, b *B.Buffer, st *HexViewState) *HexViewWidget {
	h := &HexViewWidget{id: id, buffer: b}
	h.state = st
	h.calcSizes()
	h.saveState()
	return h
}

func (h *HexViewWidget) saveState() {
}

//number of a bytes that fit in 1 line in the window width
func bytesPerLine(width, charwidth float32) int {
	//to display 1 byte takes 4 characters: 2 for hexdump, 1 trailing space and 1 print
	maxChars := int(width / (4 * charwidth))
	maxChars -= maxChars % 4
	return maxChars
}

func (h *HexViewWidget) calcSizes() {
	h.width, h.height = G.GetAvailableRegion()
	sz := I.CalcTextSize("F", true, 0)
	h.charWidth, h.charHeight = sz.X, sz.Y

	size := h.buffer.Size()
	nDigits := numHexDigits(size)
	h.addressBarWidth, _ = G.CalcTextSize(addrLabel(size, nDigits))

	h.bytesPerLine = int64(bytesPerLine(h.width-h.addressBarWidth, h.charWidth))
	h.linesPerScreen = int64(h.height / h.charHeight)

	h.topAddr = int64(I.ScrollY()/h.charHeight) * h.bytesPerLine
}

func (h *HexViewWidget) onScreen(addr int64) bool {
	top := h.topAddr
	fin := top + h.bytesPerLine*h.linesPerScreen
	return addr > top && addr < fin
}

func (h *HexViewWidget) handleKeys() {
	keymap := map[G.Key]func(){
		G.KeyJ: h.MoveDown,
		G.KeyK: h.MoveUp,
		G.KeyH: h.MoveLeft,
		G.KeyL: h.MoveRight,
	}
	for k, f := range keymap {
		if G.IsKeyPressed(k) {
			f()
		}
	}
}

func printByte(b byte) string {
	r := rune(b)
	if unicode.IsPrint(r) {
		return string(b)
	} else {
		return "."
	}
}

//TODO: use imgui columns?
//Print the cursor to the widget
func (h *HexViewWidget) printCursor(startOffset, printOffset float32) {
	cursorBG := color.RGBA{R: 255, G: 100, B: 000, A: 255}

	canvas := G.GetCanvas()
	screenPos := G.GetCursorScreenPos()
	cursor := h.state.cursor
	cursorY := int((cursor / h.bytesPerLine) * int64(h.charHeight))

	//cursor in hex dump
	cursorX := int(startOffset) + int((cursor%h.bytesPerLine)*int64(3*h.charWidth))
	pos := screenPos.Add(image.Pt(cursorX, cursorY))
	rect := image.Pt(int(h.charWidth*2), int(h.charHeight))
	canvas.AddRectFilled(pos, pos.Add(rect), cursorBG, 0, 0)

	//cursor in string dump
	cursorX = int(printOffset) + int((cursor%h.bytesPerLine)*int64(h.charWidth))
	strPos := screenPos.Add(image.Pt(cursorX, cursorY))
	strRect := image.Pt(int(h.charWidth), int(h.charHeight))
	canvas.AddRectFilled(strPos, strPos.Add(strRect), cursorBG, 0, 0)
}

func (h *HexViewWidget) inSelection(addr int64) bool {
	off, size := h.state.Selection()
	return addr >= off && addr < off+size
}

func (h *HexViewWidget) printHexBG(addr int64) {
	canvas := G.GetCanvas()
	pos := G.GetCursorScreenPos()

	//print cursor
	if addr == h.state.cursor {
		cursorBG := color.RGBA{R: 255, G: 100, B: 000, A: 255}
		rect := image.Pt(int(h.charWidth*3), int(h.charHeight))
		canvas.AddRectFilled(pos, pos.Add(rect), cursorBG, 0, 0)
	}

	//print selection
	if h.inSelection(addr) {
		selectionBG := color.RGBA{B: 200, A: 200}
		sz := 3
		if (addr+1)%h.bytesPerLine == 0 {
			sz = 2
		}
		rect := image.Pt(int(h.charWidth)*sz, int(h.charHeight))
		canvas.AddRectFilled(pos, pos.Add(rect), selectionBG, 0, 0)
	}
}

func (h *HexViewWidget) updateSelection(addr int64) {
	//adjust selection
	off := h.state.cursor
	switch {
	case addr < h.state.cursor:
		h.state.cursor = addr
		if h.state.dragging {
			h.state.selection += off - addr
		} else {
			h.state.selection += off - addr + 1
		}
	default:
		h.state.selection = addr - off + 1
	}

	//We allow EOF to be a selectable/editable field and such, chomp it here
	if h.inSelection(h.buffer.Size()) {
		h.state.selection--
	}
}

//to be called from Build() function, prints hexdump of byte and handles keyclicks and such
func (h *HexViewWidget) BuildHexCell(addr int64, b byte) {
	var hex string
	if addr == h.buffer.Size() {
		hex = "   "
	} else {
		hex = fmt.Sprintf("%02X ", b)
	}
	h.printHexBG(addr)
	I.Text(hex)

	//handle click events
	G.Event().OnClick(G.MouseButtonRight, func() {
		h.updateSelection(addr)
	}).OnMouseDown(G.MouseButtonLeft, func() {
		h.state.lastmouse = addr
		if !h.state.dragging {
			h.state.cursor = addr
			h.state.selection = 0
			h.state.dragging = true
		}
	}).OnMouseReleased(G.MouseButtonLeft, func() {
		h.state.dragging = false
		if addr != h.state.lastmouse {
			h.updateSelection(addr)
		}
	}).OnHover(func() {
		//TODO: maybe print some info i.e. interpretation of different types starting at addr
		if h.state.dragging && addr != h.state.lastmouse {
			h.updateSelection(addr)
		}
	}).Build()
}

func (h *HexViewWidget) Build() {
	I.PushStyleVarVec2(I.StyleVarFramePadding, I.Vec2{X: 0, Y: 0})
	I.PushStyleVarVec2(I.StyleVarItemSpacing, I.Vec2{X: 0, Y: 0})

	child := G.Child().Border(false).Layout(G.Custom(func() {
		h.calcSizes()
		h.handleKeys() //XXX should this be here or somewhere else??

		/* the hexview has 3 columns, so to speak,
		 * column 1 is the address, column 2 is a hexadecimal dump of bytes, column 3 is a readable string
		 * TODO: Imgui has a provision for printing text in a column. use it.
		 *
		 */
		maxAddr := numHexDigits(h.buffer.Size())                           //saved for printing address
		startOffset, _ := G.CalcTextSize(addrLabel(0, maxAddr))            //x-position of column 2
		printOffset := startOffset + h.charWidth*float32(h.bytesPerLine)*3 //x-position of column 3

		//h.printCursor(startOffset, printOffset)

		//print the hex dump using a listclipper
		numLines := (h.buffer.Size() + h.bytesPerLine - 1) / h.bytesPerLine
		lineBuffer := make([]byte, int(h.bytesPerLine)) //buffer to read the bytes for 1 line
		var clip I.ListClipper
		//dumb hack: do numlines + 10 because on big files, the last few lines get chopped off
		//due to floating point errors in scrolling calculations
		clip.BeginV(int(numLines+10), h.charHeight)
		for clip.Step() {
			for lnum := clip.DisplayStart; lnum < clip.DisplayEnd; lnum++ {
				//read data for this line
				offs := int64(lnum) * h.bytesPerLine
				h.buffer.Seek(offs, io.SeekStart)
				n, e := h.buffer.Read(lineBuffer)
				if e != nil {
					panic(e)
				}
				if n < 0 {
					break
				}

				//column 1, address
				I.Text(addrLabel(offs, maxAddr))

				//column2, hex dump
				for i := 0; i < n; i++ {
					I.SameLineV(startOffset+h.charWidth*float32(i)*3, 0)
					h.BuildHexCell(int64(lnum)*h.bytesPerLine+int64(i), lineBuffer[i])
				}
				//allow to select EOF
				if n != int(h.bytesPerLine) {
					I.SameLineV(startOffset+h.charWidth*float32(n)*3, 0)
					h.BuildHexCell(h.buffer.Size(), 0)
				}

				//column3, readable string
				for i := 0; i < n; i++ {
					if !unicode.IsGraphic(rune(lineBuffer[i])) {
						lineBuffer[i] = '.'
					}
					I.SameLineV(printOffset+float32(i)*h.charWidth, 0)
					I.Text(string(lineBuffer[i : i+1]))
				}
			}
		}
		clip.End()
	}))

	child.Build()

	I.PopStyleVarV(2) //framepadding & itemspacing
}

func (h *HexViewWidget) clampAddr(a *int64) {
	switch {
	case *a < 0:
		*a = 0
	case *a >= h.buffer.Size():
		*a = h.buffer.Size() - 1
	}
}

func (h *HexViewWidget) MoveRight() {
	h.state.cursor += 1
	h.finishMove()
}

func (h *HexViewWidget) MoveLeft() {
	h.state.cursor -= 1
	h.finishMove()
}

func (h *HexViewWidget) MoveDown() {
	h.state.cursor += h.bytesPerLine
	h.finishMove()
}

func (h *HexViewWidget) MoveUp() {
	h.state.cursor -= h.bytesPerLine
	h.finishMove()
}

func (h *HexViewWidget) finishMove() {
	h.state.selection = 0
	h.clampAddr(&h.state.cursor)
	h.ScrollTo(h.state.cursor)
	h.saveState()
}

func (h *HexViewWidget) ScrollTo(addr int64) {
	bpl := h.bytesPerLine
	top := h.topAddr
	switch {
	case h.state.cursor < top:
		//scroll up, addr should be in the first line
		//make first line the one that contains addr
		top = (h.state.cursor / bpl) * bpl

	case h.state.cursor > top+bpl*h.linesPerScreen:
		//scroll down, addr should be in the last line
		a := h.state.cursor - h.linesPerScreen*bpl
		top = ((a + bpl - 1) / bpl) * bpl
	default:
		//addr is already on screen
	}
	h.clampAddr(&top)
	I.SetScrollY(float32(top/bpl) * h.charHeight)
	h.saveState()
}
