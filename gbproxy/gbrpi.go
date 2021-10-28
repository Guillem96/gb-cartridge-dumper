package gbproxy

import (
	"github.com/Guillem96/gb-dumper/bytes"
	"github.com/Guillem96/gb-dumper/io"
	"github.com/stianeikeland/go-rpio/v4"
)

// GameBoyRPiPin implements the GameBoyPin interface. This implementation maps connections between
// a RaspberryPi and the GameBoy via GPIO
type GameBoyRPiPin rpio.Pin

// Read returns the RaspberryPi GPIO pin status
func (p GameBoyRPiPin) Read() bool {
	return rpio.ReadPin(rpio.Pin(p)) == rpio.High
}

// High sets the the RaspberryPi GPIO pin status to high
func (p GameBoyRPiPin) High() {
	rpio.WritePin(rpio.Pin(p), rpio.High)
}

// Low sets the the RaspberryPi GPIO pin status to low
func (p GameBoyRPiPin) Low() {
	rpio.WritePin(rpio.Pin(p), rpio.Low)
}

// SetState sets the the RaspberryPi GPIO pin to the given status
func (p GameBoyRPiPin) SetState(state bool) {
	if state {
		rpio.Pin(p).High()
	} else {
		rpio.Pin(p).Low()
	}
}

// Input sets the the RaspberryPi GPIO pin mode to input (populated by the GameBoy)
func (p GameBoyRPiPin) Input() {
	rpio.PinMode(rpio.Pin(p), rpio.Input)
}

// Output sets the the RaspberryPi GPIO pin mode to input (populated by the RPi)
func (p GameBoyRPiPin) Output() {
	rpio.PinMode(rpio.Pin(p), rpio.Output)
}

// RPiGameBoyProxy implements the GameBoyProxy to provide a working data transfer between
// a RaspberryPi and the GameBoy
type RPiGameBoyProxy struct {
	As []GameBoyRPiPin
	Db []GameBoyRPiPin
	Rd GameBoyRPiPin
	Wr GameBoyRPiPin
}

// NewRPiGameBoyProxy creates a new RPiGameBoyProxy
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

	rd := GameBoyRPiPin(cm.RD)
	rd.Output()

	wr := GameBoyRPiPin(cm.WR)
	wr.Output()

	return &RPiGameBoyProxy{
		As: as,
		Db: db,
		Rd: rd,
		Wr: wr,
	}
}

// Read reads the byte located in the address specified with the SelectAddress method.
func (rpigb *RPiGameBoyProxy) Read() uint8 {
	var result uint8

	// When reading we set DX pins to input mode
	for _, d := range rpigb.Db {
		d.Input()
	}

	// TODO: Set RD pin high?

	result = 0x00
	for i := 0; i < 8; i++ {
		if rpigb.Db[i].Read() {
			result += (1 << i)
		}
	}
	return result
}

// Write writes the provided value to the selected address with the SelectAddress function
func (rpigb *RPiGameBoyProxy) Write(value uint) error {
	// When writing we set DX pins to output mode
	for _, d := range rpigb.Db {
		d.Output()
	}

	pinsState, err := bytes.AddressToBitArray(value, 8)
	if err != nil {
		return err
	}

	for i, ps := range pinsState {
		if ps {
			rpigb.Db[i].High()
		} else {
			rpigb.Db[i].Low()
		}
	}

	// TODO: Set WR pin high?
	return nil
}

// SelectAddress sets the GPIO pins status so the referenced address in the cartridge is the given one
func (rpigb *RPiGameBoyProxy) SelectAddress(addr uint) error {
	pinsState, err := bytes.AddressToBitArray(addr, 16)
	if err != nil {
		return err
	}

	for i, ps := range pinsState {
		if ps {
			rpigb.As[i].High()
		} else {
			rpigb.As[i].Low()
		}
	}
	return nil
}
