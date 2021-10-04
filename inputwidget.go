package main

import (
	"fmt"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
)

//simple widget that allows entering 2 hexadecimal characters

type InputHexByte struct {
	id        string
	hi, lo    byte //should be <16
	at        int  //0 or 1 (hi or low nibble)
	cbCancel  func()
	cbSuccess func(b byte)
}

func (ih *InputHexByte) Dispose() {
	//empty
}

func InputHex(id string, cancel func(), success func(b byte)) G.Widget {
	hRaw := G.Context.GetState(id)
	var h *InputHexByte
	if hRaw != nil {
		h = hRaw.(*InputHexByte)
		h.cbCancel = cancel
		h.cbSuccess = success
	} else {
		h = &InputHexByte{
			cbCancel:  cancel,
			cbSuccess: success,
		}
	}
	G.Context.SetState(id, h)
	return h
}

func (ih *InputHexByte) Build() {
	var txt string
	switch ih.at {
	case 0:
		txt = "__ "
	case 1:
		txt = fmt.Sprintf("%1X_ ", ih.hi)
	case 2:
		txt = fmt.Sprintf("%1X%1X ", ih.hi, ih.lo)
	}

	I.Text(txt)

	var input byte = 255
	switch {
	case G.IsKeyPressed(G.Key0):
		input = 0
	case G.IsKeyPressed(G.Key1):
		input = 1
	case G.IsKeyPressed(G.Key2):
		input = 2
	case G.IsKeyPressed(G.Key3):
		input = 3
	case G.IsKeyPressed(G.Key4):
		input = 4
	case G.IsKeyPressed(G.Key5):
		input = 5
	case G.IsKeyPressed(G.Key6):
		input = 6
	case G.IsKeyPressed(G.Key7):
		input = 7
	case G.IsKeyPressed(G.Key8):
		input = 8
	case G.IsKeyPressed(G.Key9):
		input = 9
	case G.IsKeyPressed(G.KeyA):
		input = 0xA
	case G.IsKeyPressed(G.KeyB):
		input = 0xB
	case G.IsKeyPressed(G.KeyC):
		input = 0xC
	case G.IsKeyPressed(G.KeyD):
		input = 0xD
	case G.IsKeyPressed(G.KeyE):
		input = 0xE
	case G.IsKeyPressed(G.KeyF):
		input = 0xF
	case G.IsKeyPressed(G.KeyBackspace):
		ih.at--
		if ih.at < 0 {
			ih.cancel()
			return
		}
	case G.IsKeyPressed(G.KeyEscape):
		ih.cancel()
		return
	}
	if input != 255 {
		if ih.at == 0 {
			ih.hi = input
		} else {
			ih.lo = input
		}
		ih.at++
		if ih.at >= 2 {
			ih.success()
		}
	}
	G.Context.SetState(ih.id, ih)
}

func (ih *InputHexByte) cancel() {
	if ih.cbCancel != nil {
		ih.cbCancel()
	}
	ih.reset()
}

func (ih *InputHexByte) reset() {
	ih.hi = 0
	ih.lo = 0
	ih.at = 0
}

func (ih *InputHexByte) success() {
	if ih.cbSuccess != nil {
		ih.cbSuccess(ih.hi<<4 | ih.lo)
		ih.reset()
	}
}
