package cartridge

import "github.com/Guillem96/gb-dumper/gbproxy"

type Dumper struct {
	gbp gbproxy.GameBoyProxy
}

func NewDumper(gbp gbproxy.GameBoyProxy) *Dumper {
	return &Dumper{
		gbp: gbp,
	}
}

// Reads the stored byte in cartridge requested address
func (d *Dumper) Read(addr uint16) uint8 {
	d.gbp.SelectAddress(addr)

	// TODO: Should we set the RD pin to Low?
	// TODO: Should we sleep here?

	return d.gbp.Read()
}

// Reads a range of addresses
func (d *Dumper) ReadRange(startAddr uint16, endAddr uint16) []uint8 {
	return []uint8{0x00, 0x00}
}

func (d *Dumper) ReadHeader() *CartridgeHeader {
	bytes := d.ReadRange(0x00, 0x150)
	return ROMHeaderFromBytes(bytes)
}
