package gbproxy

import (
	"github.com/Guillem96/gb-dumper/bytes"
	"github.com/Guillem96/gb-dumper/io"
	"github.com/stianeikeland/go-rpio"
)

type GameBoyRPiPin rpio.Pin

func (p GameBoyRPiPin) Read() bool {
	return rpio.ReadPin(rpio.Pin(p)) == rpio.High
}

func (p GameBoyRPiPin) High() {
	rpio.WritePin(rpio.Pin(p), rpio.High)
}

func (p GameBoyRPiPin) Low() {
	rpio.WritePin(rpio.Pin(p), rpio.Low)
}

func (p GameBoyRPiPin) SetState(state bool) {
	if state {
		rpio.Pin(p).High()
	} else {
		rpio.Pin(p).Low()
	}
}

func (p GameBoyRPiPin) Input() {
	rpio.PinMode(rpio.Pin(p), rpio.Input)
}

func (p GameBoyRPiPin) Output() {
	rpio.PinMode(rpio.Pin(p), rpio.Output)
}

type RPiGameBoyProxy struct {
	As []GameBoyRPiPin
	Db []GameBoyRPiPin
}

func NewRPiGameBoyProxy(cm *io.GameBoyRaspberryMapping) *RPiGameBoyProxy {
	as := []GameBoyRPiPin{GameBoyRPiPin(cm.A0), GameBoyRPiPin(cm.A1), GameBoyRPiPin(cm.A2),
		GameBoyRPiPin(cm.A3), GameBoyRPiPin(cm.A4), GameBoyRPiPin(cm.A5), GameBoyRPiPin(cm.A6),
		GameBoyRPiPin(cm.A7), GameBoyRPiPin(cm.A8), GameBoyRPiPin(cm.A9), GameBoyRPiPin(cm.A10),
		GameBoyRPiPin(cm.A11), GameBoyRPiPin(cm.A12), GameBoyRPiPin(cm.A13), GameBoyRPiPin(cm.A14),
		GameBoyRPiPin(cm.A15),
	}

	db := []GameBoyRPiPin{GameBoyRPiPin(cm.D0), GameBoyRPiPin(cm.D1), GameBoyRPiPin(cm.D2),
		GameBoyRPiPin(cm.D3), GameBoyRPiPin(cm.D4), GameBoyRPiPin(cm.D5), GameBoyRPiPin(cm.D6),
		GameBoyRPiPin(cm.D7),
	}

	// AX pins are the address selector.
	for _, a := range as {
		a.Output()
	}

	// We read the data stored in the requested addres (DX pins)
	for _, d := range db {
		d.Input()
	}

	return &RPiGameBoyProxy{
		As: as,
		Db: db,
	}
}

func (rpigb *RPiGameBoyProxy) Read() uint8 {
	var result uint8
	result = 0x00
	for i := 0; i < 8; i++ {
		if rpigb.Db[i].Read() {
			result += (1 << i)
		}
	}
	return result
}

func (rpigb *RPiGameBoyProxy) SelectAddress(addr uint16) {
	pinsState := bytes.AddressToBitArray(addr)
	for i, ps := range pinsState {
		if ps {
			rpigb.As[i].High()
		} else {
			rpigb.As[i].Low()
		}
	}
}
