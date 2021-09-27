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
	addr int64 //from start of file; byte address of the cursor
	//which index of string rep of the data on this addr is the cursor?
	idx int //only for insert mode
}

//Widget state
type HexViewState struct {
	topAddr int64
	cursor  cursorAddr
}

//implements: giu.Disposable
func (st *HexViewState) Dispose() {
	//nothing to do here
}

type HexViewWidget struct {
	state *HexViewState

	id     string
	buffer *B.Buffer

	bytesPerLine    int64
	linesPerScreen  int64
	width           float32
	height          float32
	charWidth       float32
	charHeight      float32
	addressBarWidth float32
}

func HexView(id string, b *B.Buffer) *HexViewWidget {
	h := &HexViewWidget{id: id, buffer: b}
	st, ok := G.Context.GetState(id).(*HexViewState)
	if ok {
		h.state = st
	} else {
		h.state = &HexViewState{}
	}
	h.calcSizes()
	h.saveState()
	return h
}

func (h *HexViewWidget) saveState() {
	G.Context.SetState(h.id, h.state)
}

//number of a bytes that fit in 1 line in the window width
func bytesPerLine(width, charwidth float32) int {
	//to display 1 byte takes 4 characters: 2 for hexdump, 1 trailing space and 1 print
	maxChars := int(width / (4 * charwidth))
	return maxChars //TODO: round down to multiple of 4, 8 or 16?
}

func (h *HexViewWidget) calcSizes() {
	h.width, h.height = G.GetAvailableRegion()
	//h.charWidth, h.charHeight = G.CalcTextSize("F")
	//XXX somehow charWidth is off by a factor of 2...?????
	h.charWidth, h.charHeight = G.CalcTextSize("")

	size := h.buffer.Size()
	nDigits := numHexDigits(size)
	h.addressBarWidth, _ = G.CalcTextSize(addrLabel(size, nDigits))

	h.bytesPerLine = int64(bytesPerLine(h.width-h.addressBarWidth, h.charWidth))
	h.linesPerScreen = int64(h.height / h.charHeight)
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

func (h *HexViewWidget) Build() {
	cursorBG := color.RGBA{R: 255, G: 100, B: 000, A: 255}
	I.PushStyleVarVec2(I.StyleVarFramePadding, I.Vec2{X: 0, Y: 0})
	I.PushStyleVarVec2(I.StyleVarItemSpacing, I.Vec2{X: 0, Y: 0})

	child := G.Child().Layout(G.Custom(func() {
		h.calcSizes()
		h.handleKeys() //XXX should this be here or somewhere else??

		/* the hexview has 3 columns, so to speak,
		 * column 1 is the address, column 2 is a hexadecimal dump of bytes, column 3 is a readable string
		 */
		maxAddr := numHexDigits(h.buffer.Size())                           //saved for printing address
		startOffset, _ := G.CalcTextSize(addrLabel(0, maxAddr))            //x-position of column 2
		printOffset := startOffset + h.charWidth*float32(h.bytesPerLine)*3 //x-position of column 3

		numLines := (h.buffer.Size() + h.bytesPerLine - 1) / h.bytesPerLine //round up
		lineBuffer := make([]byte, int(h.bytesPerLine))                     //buffer to read the bytes for 1 line

		var clip I.ListClipper
		clip.Begin(int(numLines))
		for clip.Step() {
			for lnum := clip.DisplayStart; lnum < clip.DisplayEnd; lnum++ {
				offs := int64(lnum) * h.bytesPerLine
				h.buffer.Seek(offs, io.SeekStart)
				n, _ := h.buffer.Read(lineBuffer)
				fmt.Printf("line %d/%d @ %X, read %d bytes, len(lineBuf) = %d\n", lnum, numLines, offs, n, len(lineBuffer))

				//print address (column 1)
				I.Text(addrLabel(offs, maxAddr))

				//print hexdump (column 2)
				for i := 0; i < n; i++ {
					I.SameLineV(startOffset+h.charWidth*float32(i)*3, 0)
					addr := offs + int64(i)
					if addr == h.state.cursor.addr {
						canvas := G.GetCanvas()

						//cursor shows in hex dump
						pos := G.GetCursorScreenPos()
						rect := image.Point{int(h.charWidth * 3), int(h.charHeight)}
						canvas.AddRectFilled(pos, pos.Add(rect), cursorBG, 0, 0)

						//cursor shows in string dump
						strPos := image.Point{X: 1 + int(printOffset) + int(h.charWidth)*(i+1), Y: pos.Y}
						strRect := image.Point{X: int(h.charWidth), Y: int(h.charHeight)}
						canvas.AddRectFilled(strPos, strPos.Add(strRect), cursorBG, 0, 0)
					}
					hex := fmt.Sprintf("%02X ", lineBuffer[i])
					I.Text(hex)
				}
				//print readable string
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
	h.state.cursor.addr += 1
	h.state.cursor.idx = 0
	h.finishMove()
}

func (h *HexViewWidget) MoveLeft() {
	h.state.cursor.addr -= 1
	h.state.cursor.idx = 0
	h.finishMove()
}

func (h *HexViewWidget) MoveDown() {
	h.state.cursor.addr += h.bytesPerLine
	h.state.cursor.idx = 0
	h.finishMove()
}

func (h *HexViewWidget) MoveUp() {
	h.state.cursor.addr -= h.bytesPerLine
	h.state.cursor.idx = 0
	h.finishMove()
}

func (h *HexViewWidget) finishMove() {
	h.clampAddr(&h.state.cursor.addr)
	/* TODO
	if h.state.cursor.addr is outside view {
		h.ScrollTo(h.state.cursor.addr)
	}
	*/
	h.saveState()
}

func (h *HexViewWidget) ScrollTo(addr int64) {
	bpl := h.bytesPerLine
	h.state.cursor.addr = addr
	h.clampAddr(&h.state.cursor.addr)
	switch {
	case h.state.cursor.addr < h.state.topAddr:
		//scroll up, addr should be in the first line
		//make first line the one that contains addr
		h.state.topAddr = (h.state.cursor.addr / bpl) * bpl

	case h.state.cursor.addr > h.state.topAddr+bpl*h.linesPerScreen:
		//scroll down, addr should be in the last line
		a := h.state.cursor.addr - h.linesPerScreen*bpl
		h.state.topAddr = ((a + bpl - 1) / bpl) * bpl
	default:
		//addr is already on screen
	}
	h.clampAddr(&h.state.topAddr)
	h.saveState()
}
