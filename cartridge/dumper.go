package cartridge

import (
	"github.com/Guillem96/gb-dumper/io"
	"github.com/stianeikeland/go-rpio/v4"
)

type Dumper struct {
	addressSelector []rpio.Pin
	dataBus         []rpio.Pin
}

func NewDumper(cm *io.GameBoyRaspberryMapping) *Dumper {
	as := []rpio.Pin{rpio.Pin(cm.A0), rpio.Pin(cm.A1), rpio.Pin(cm.A2), rpio.Pin(cm.A3), rpio.Pin(cm.A4),
		rpio.Pin(cm.A5), rpio.Pin(cm.A6), rpio.Pin(cm.A7), rpio.Pin(cm.A8), rpio.Pin(cm.A9),
		rpio.Pin(cm.A10), rpio.Pin(cm.A11), rpio.Pin(cm.A12), rpio.Pin(cm.A13), rpio.Pin(cm.A14),
		rpio.Pin(cm.A15),
	}

	db := []rpio.Pin{rpio.Pin(cm.D0), rpio.Pin(cm.D1), rpio.Pin(cm.D2), rpio.Pin(cm.D3), rpio.Pin(cm.D4),
		rpio.Pin(cm.D5), rpio.Pin(cm.D6), rpio.Pin(cm.D7),
	}

	// AX pins are the address selector.
	for _, a := range as {
		a.Output()
	}

	// We read the data stored in the requested addres (DX pins)
	for _, d := range db {
		d.Input()
	}

	return &Dumper{
		addressSelector: as,
		dataBus:         db,
	}
}

// Reads the stored byte in cartridge requested address
func (d *Dumper) ReadAddress(addr uint16) uint8 {
	pinsState := addrToPinsState(addr)
	for i, pinSt := range pinsState {
		rpio.WritePin(d.addressSelector[i], pinSt)
	}

	// TODO: Should we set the RD pin to Low?
	// TODO: Should we sleep here?

	return d.readCurrentData()
}

// Reads a range of addresses
func (d *Dumper) ReadAddressRange(startAddr uint16, endAddr uint16) []uint8 {
	return []uint8{0x00, 0x00}
}

func (d *Dumper) readCurrentData() (result uint8) {
	result = 0x00
	for pos, pin := range d.dataBus {
		if pin.Read() == rpio.High {
			result += (1 << pos)
		}
	}
	return result
}

func addrToPinsState(addr uint16) []rpio.State {
	pins := make([]rpio.State, 18)
	for i := 0; i < 16; i++ {
		if testBit(addr, uint(i)) {
			pins[i] = rpio.High
		} else {
			pins[i] = rpio.Low
		}
	}
	return pins
}

func testBit(n uint16, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}
