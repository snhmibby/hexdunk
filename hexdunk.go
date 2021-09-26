package main

//TODO: most hex dumps have a printable character view after the raw data, implement it
//TODO: use timestamps to check to see if file is edited by another program

import (
	"fmt"
	"image"
	"image/color"
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
	buffer *B.FileBuffer

	bytesPerLine    int64
	linesPerScreen  int64
	width           float32
	height          float32
	charWidth       float32
	charHeight      float32
	addressBarWidth float32
}

func HexView(id string, b *B.FileBuffer) *HexViewWidget {
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
	return maxChars //TODO: round down to multiple of 8 or 16?
}

func (h *HexViewWidget) calcSizes() {
	h.width, h.height = G.GetAvailableRegion()
	//h.charWidth, h.charHeight = G.CalcTextSize("F")
	//XXX somehow our calculation is off by a factor of 2 but i have no clue how...
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
	I.PushStyleVarVec2(I.StyleVarFramePadding, I.Vec2{X: 0, Y: 0})
	I.PushStyleVarVec2(I.StyleVarItemSpacing, I.Vec2{X: 0, Y: 0})

	h.calcSizes()
	h.handleKeys() //should this be here?

	buf, _ := h.buffer.ReadBuf(h.state.topAddr, h.bytesPerLine*h.linesPerScreen)
	var found_eof = false
	maxAddr := numHexDigits(h.buffer.Size()) //cached for printing
	for lnum := int64(0); lnum < h.linesPerScreen; lnum++ {
		if found_eof {
			break
		}
		lineStr := ""
		lineAddr := h.state.topAddr + lnum*h.bytesPerLine
		//TODO: colors? color addr if cursor on this line?
		I.Text(addrLabel(lineAddr, maxAddr))

		//hexdump 1 line
		for addrOffset := int64(0); addrOffset < h.bytesPerLine; addrOffset++ {
			I.SameLine()
			b, err := buf.ReadByte()
			if err != nil {
				found_eof = true
				I.Text("   ")
			} else {
				lineStr += printByte(b)
				hex := fmt.Sprintf("%02X ", b)
				rect := image.Point{int(h.charWidth * 3), int(h.charHeight)}
				isCursor := lineAddr+addrOffset == h.state.cursor.addr
				if isCursor {
					//just a red background for cursor for now
					cursorBG := color.RGBA{R: 255, G: 100, B: 000, A: 255}
					pos := G.GetCursorScreenPos()
					can := G.GetCanvas()
					can.AddRectFilled(pos, pos.Add(rect), cursorBG, 0, 0)
				}
				I.Text(hex)
			}
		}
		I.SameLine()
		I.Text(lineStr)
	}
	I.PopStyleVarV(2) //framepadding & itemspacing
}

func (h *HexViewWidget) clampAddr(a *int64) {
	switch {
	case *a < 0:
		*a = 0
	case *a > h.buffer.Size():
		*a = h.buffer.Size()
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

//for testing

var fileName string
var openBuffer *B.FileBuffer

func draw() {
	hv := HexView("hexview##1", openBuffer)
	G.SingleWindow().Layout(
		G.Label(fileName),
		G.Child().Layout(hv),
	)
}

func main() {
	var err error
	fileName = "go.mod"
	openBuffer, err = B.NewFileBuffer(fileName)
	if err != nil {
		panic(err)
	}
	G.SetDefaultFont("DejavuSansMono.ttf", 11)
	w := G.NewMasterWindow("HexViewer", 800, 600, 0)
	w.Run(draw)
}
