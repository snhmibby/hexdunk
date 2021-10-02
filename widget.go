package main

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
	return fmt.Sprintf("%0*X:", nDigits, addr)
}

type TextColor struct {
	fg, bg color.RGBA
}

type cursorAddr struct {
	//byte address of cursor
	addr int64
}

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

func bytesPerLine(width, charwidth float32) int {
	//to display 1 byte takes 4 characters: 2 for hexdump, 1 trailing space and 1 print
	maxChars := int(width / (4 * charwidth))

	//round to multiple of 4
	maxChars -= maxChars % 4
	if maxChars == 0 {
		maxChars = 4
	}
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
	if unicode.IsGraphic(rune(b)) {
		return string(b)
	} else {
		return "."
	}
}
func (h *HexViewWidget) inSelection(addr int64) bool {
	off, size := h.state.Selection()
	return addr >= off && addr < off+size
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

func (h *HexViewWidget) printBG(addr int64, cursorw, selectw int) {
	canvas := G.GetCanvas()
	pos := G.GetCursorScreenPos()

	if addr == h.state.cursor {
		cursorBG := color.RGBA{R: 255, G: 100, B: 000, A: 255}
		rect := image.Pt(cursorw*int(h.charWidth), int(h.charHeight))
		canvas.AddRectFilled(pos, pos.Add(rect), cursorBG, 0, 0)
	}

	if h.inSelection(addr) {
		selectionBG := color.RGBA{B: 200, A: 100}
		rect := image.Pt(selectw*int(h.charWidth), int(h.charHeight))
		canvas.AddRectFilled(pos, pos.Add(rect), selectionBG, 0, 0)
	}
}

//a cell is a piece of text that corresponds to an file-offset.
//it can be clicked and dragged
func (h *HexViewWidget) BuildCell(addr int64, txt string) {
	if len(txt) == 1 {
		h.printBG(addr, 1, 1)
	} else {
		h.printBG(addr, 2, 3)
	}
	I.Text(txt)

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

//to be called from Build() function, prints hexdump of byte and handles keyclicks and such
func (h *HexViewWidget) BuildHexCell(addr int64, b byte) {
	var hex string
	if addr == h.buffer.Size() {
		hex = "   "
	} else {
		hex = fmt.Sprintf("%02X ", b)
	}
	h.BuildCell(addr, hex)
}

//to be called from Build() function, prints readable interpretation of byte and handles keyclicks and such
func (h *HexViewWidget) BuildStrCell(addr int64, b byte) {
	str := printByte(b)
	h.BuildCell(addr, str)
}

func (h *HexViewWidget) Build() {
	I.PushStyleVarVec2(I.StyleVarFramePadding, I.Vec2{})
	I.PushStyleVarVec2(I.StyleVarItemSpacing, I.Vec2{})
	I.PushStyleVarVec2(I.StyleVarCellPadding, I.Vec2{})
	defer I.PopStyleVarV(3)

	child := G.Child().Border(false).Layout(G.Custom(func() {
		h.calcSizes()
		h.handleKeys() //XXX should this be here or somewhere else??

		maxAddr := numHexDigits(h.buffer.Size())                //saved for printing address
		addressSize, _ := G.CalcTextSize(addrLabel(0, maxAddr)) //x-position of column 2

		flags := I.TableFlags_BordersInnerV | I.TableFlags_BordersOuter | I.TableFlags_SizingFixedFit
		if I.BeginTable("HexDumpTable", 3, flags, I.Vec2{}, 0) {
			defer I.EndTable()
			I.TableSetupColumn("Offset", 0, addressSize, 0)
			I.TableSetupColumn("HexDump", 0, 3*h.charWidth*float32(h.bytesPerLine), 0)
			I.TableSetupColumn("Readable", 0, float32(h.bytesPerLine)*h.charWidth, 0)

			//print the hex dump using a listclipper
			numLines := (h.buffer.Size() + h.bytesPerLine - 1) / h.bytesPerLine
			lineBuffer := make([]byte, int(h.bytesPerLine)) //buffer to read the bytes for 1 line
			var clip I.ListClipper
			//dumb hack: do numlines + 10 because on big files, the last few lines get chopped off
			//due to floating point errors in scrolling calculations
			clip.BeginV(int(numLines+10), h.charHeight)
			defer clip.End()
			for clip.Step() {
				for lnum := clip.DisplayStart; lnum < clip.DisplayEnd; lnum++ {
					offs := int64(lnum) * h.bytesPerLine
					if offs > h.buffer.Size() {
						break
					}

					//read data for this line
					h.buffer.Seek(offs, io.SeekStart)
					n, e := h.buffer.Read(lineBuffer)
					if n < 0 || e != nil {
						//TODO: properly handle EOF here!!!
						break
					}

					//address
					I.TableNextColumn()
					I.Text(addrLabel(offs, maxAddr))

					//hex dump
					I.TableNextColumn()
					for i := 0; i < n; i++ {
						if i != 0 {
							I.SameLine()
						}
						h.BuildHexCell(offs+int64(i), lineBuffer[i])
					}
					//allow to select EOF??
					if n != int(h.bytesPerLine) {
						if n != 0 {
							I.SameLine()
						}
						h.BuildHexCell(h.buffer.Size(), 0)
					}

					//readable string
					I.TableNextColumn()
					for i := 0; i < n; i++ {
						if i != 0 {
							I.SameLine()
						}
						h.BuildStrCell(offs+int64(i), lineBuffer[i])
					}
				}
			}
		}
	}))

	child.Build()

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
