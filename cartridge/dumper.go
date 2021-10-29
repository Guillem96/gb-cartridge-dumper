package cartridge

import (
	"errors"
	"log"
	"os"

	"github.com/Guillem96/gb-dumper/gbproxy"
)

// Dumper is the object responsible of running "queries" against the cartridge dumper
// The queries are send using the provided gbproxy.GameBoyProxy
type Dumper struct {
	l      *log.Logger
	gbp    gbproxy.GameBoyProxy
	header *CartridgeHeader
}

// NewDumper creates a new cartridge dumper and returns a pointer to it
func NewDumper(gbp gbproxy.GameBoyProxy) *Dumper {
	return &Dumper{
		gbp:    gbp,
		header: nil,
		l:      log.New(os.Stdout, "[GB Cartridge Dumper]", log.LstdFlags),
	}
}

// Read reads the stored byte in cartridge requested address
func (d *Dumper) Read(addr uint) uint8 {
	d.gbp.SelectAddress(addr)
	return d.gbp.Read()
}

// ReadRange reads a range of addresses and returns all read bytes in order
func (d *Dumper) ReadRange(startAddr uint, endAddr uint) []uint8 {
	d.l.Printf("Reading address range 0x%x to 0x%x", startAddr, endAddr)
	rb := make([]uint8, 0)
	for ca := startAddr; ca < endAddr; ca++ {
		rb = append(rb, d.Read(ca))
	}
	return rb
}

// ReadHeader reads the whole cartridge header
func (d *Dumper) ReadHeader() *CartridgeHeader {
	if d.header != nil {
		return d.header
	}
	d.l.Println("Reading ROM header data.")
	bytes := d.ReadRange(0x00, 0x150)
	d.header = ROMHeaderFromBytes(bytes)
	return d.header
}

// ChangeROMBank communicates with the GameBoy MBC and changes the active ROM bank
func (d *Dumper) ChangeROMBank(bank uint) error {
	h := d.ReadHeader()

	if !h.HasMBC() {
		return errors.New("cartridge has no MBC.")
	}

	if h.IsMBC1() {
		if bank > 127 {
			return errors.New("MBC1 only has 7 bits to address a ROM bank. The provided bank cannot be represented using 7bits")
		}

		// Select ROM banking mode
		d.gbp.SelectAddress(0x6000)
		d.gbp.Write(uint8(1))

		// Low bank number
		d.gbp.SelectAddress(0x2100)
		d.gbp.Write(uint8(bank & 0x1F))

		// 2 bit high number
		d.gbp.SelectAddress(0x4000)
		d.gbp.Write(uint8((bank >> 5) & 0x03))
	} else {
		return errors.New("cartridge not supported yet.")
	}

	return nil
}
