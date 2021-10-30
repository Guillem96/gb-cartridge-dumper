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

// ReadCartridge dumps the whole cartridge data. Reads all ROM & RAM banks
func (d *Dumper) ReadCartridge() *Cartridge {
	h := d.ReadHeader()

	// Dump Rom banks
	nb := h.GetNumROMBanks()
	banks := make([][]uint8, nb)
	d.l.Printf("The cartridge has %d banks.\n", nb)

	var addrBase uint
	addrBase = 0x0000
	for b := 0; b < nb; b++ {
		if h.IsMBC1() && (nb == 0x00 || nb == 0x20 || nb == 0x40 || nb == 0x60) {
			addrBase = 0x0000
		} else if b > 0 {
			addrBase = 0x4000
		}

		banks[b] = d.ReadRange(addrBase, addrBase+0x4000)
	}

	// TODO: Dump the Cartridge RAM

	return NewCartridge(h, banks)
}

// ChangeROMBank communicates with the GameBoy MBC and changes the active ROM bank
// MBC1 cartridges map the bank 0x20 0x40 and 0x60 to 0x0000-0x3FFF address (caller must be aware of this)
func (d *Dumper) ChangeROMBank(bank uint) error {
	h := d.ReadHeader()
	nb := h.GetNumROMBanks()

	if !h.HasMBC() {
		return errors.New("cartridge has no MBC.")
	}

	if bank > uint(nb-1) {
		errMsg := fmt.Sprintf("cartridge type 0x%x only has %d banks. You want to change to bank %d.\n",
			h.CartridgeType, nb, bank)
		return errors.New(errMsg)
	}

	// Compute the low bit mask
	var lbm uint
	if h.IsMBC1() {
		lbm = 0x1F
	} else if h.IsMBC2() {
		lbm = 0x0F
	} else if h.IsMBC3() {
		lbm = 0x7F
	} else if h.IsMBC5() {
		lbm = 0xFF
	} else {
		return errors.New("cartridge type not supported yet.")
	}

	if h.IsMBC1() {
		// Enable ROM banking mode (depending on the bank we select advanced or normal banking)
		d.gbp.SelectAddress(0x6000)

		lb := bank & lbm
		if lb == 0 {
			// If all lower bank bits are 0 enable advanced ROM banking so we can map
			// the selected bank to 0x0000 to 0x3FFF
			d.gbp.Write(uint8(1))
		} else {
			d.gbp.Write(uint8(0))
		}

		// 2 bit high number for MBC1
		d.gbp.SelectAddress(0x4000)
		d.gbp.Write(uint8((bank >> 5) & 0x03))
	} else if h.IsMBC5() {
		// 1 bit high number for MBC5
		d.gbp.SelectAddress(0x4000)
		d.gbp.Write(uint8((bank >> 8) & 0x01))
	}

	// Low bank number
	d.gbp.SelectAddress(0x2100)
	d.gbp.Write(uint8(bank & lbm))

	return nil
}
