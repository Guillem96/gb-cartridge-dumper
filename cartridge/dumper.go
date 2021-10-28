package cartridge

import (
	"errors"
	"fmt"
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
func (d *Dumper) Read(addr uint) (uint8, error) {
	if err := d.gbp.SelectAddress(addr); err != nil {
		return 0x0, err
	}
	return d.gbp.Read(), nil
}

// ReadRange reads a range of addresses and returns all read bytes in order
func (d *Dumper) ReadRange(startAddr uint, endAddr uint) ([]uint8, error) {
	d.l.Printf("Reading address range 0x%x to 0x%x", startAddr, endAddr)
	rb := make([]uint8, 0)
	for ca := startAddr; ca < endAddr; ca++ {
		v, err := d.Read(ca)
		if err != nil {
			errMsg := fmt.Sprintf("reading address 0x%x: %v", ca, err)
			return nil, errors.New(errMsg)
		}
		rb = append(rb, v)
	}
	return rb, nil
}

// ReadHeader reads the whole cartridge header
func (d *Dumper) ReadHeader() (*CartridgeHeader, error) {
	d.l.Println("Reading ROM header data.")
	bytes, err := d.ReadRange(0x00, 0x150)
	if err != nil {
		return nil, err
	}
	d.header = ROMHeaderFromBytes(bytes)
	return d.header, nil
}
