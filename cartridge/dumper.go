package cartridge

import "github.com/Guillem96/gb-dumper/gbproxy"

// Dumper is the object responsible of running "queries" against the cartridge dumper
// The queries are send using the provided gbproxy.GameBoyProxy
type Dumper struct {
	gbp gbproxy.GameBoyProxy
}

// NewDumper creates a new cartridge dumper and returns a pointer to it
func NewDumper(gbp gbproxy.GameBoyProxy) *Dumper {
	return &Dumper{
		gbp: gbp,
	}
}

// Read reads the stored byte in cartridge requested address
func (d *Dumper) Read(addr uint16) uint8 {
	d.gbp.SelectAddress(addr)

	// TODO: Should we set the RD pin to Low?
	// TODO: Should we sleep here?

	return d.gbp.Read()
}

// ReadRange reads a range of addresses and returns all read bytes in order
func (d *Dumper) ReadRange(startAddr uint16, endAddr uint16) []uint8 {
	rb := make([]uint8, 0)
	for ca := startAddr; ca < endAddr; ca++ {
		rb = append(rb, d.Read(ca))
	}
	return rb
}

// ReadHeader reads the whole cartridge header
func (d *Dumper) ReadHeader() *CartridgeHeader {
	bytes := d.ReadRange(0x00, 0x150)
	return ROMHeaderFromBytes(bytes)
}
