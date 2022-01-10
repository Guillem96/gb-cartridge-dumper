package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Guillem96/gameboy-tools/cartridge"
	"github.com/Guillem96/gameboy-tools/conmap"
	"github.com/Guillem96/gameboy-tools/gbproxy"
)

func printByteArray(arr []uint8, rows int) {
	rpr := int(math.Ceil(float64(len(arr)) / float64(rows)))
	for i, b := range arr {
		if i%rpr == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("%02x ", b)
	}
	fmt.Println()
}

func Equal(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type DifferRecord struct {
	valueA uint8
	valueB uint8
	pos    int
}

type DifferRecords []DifferRecord

func (ds DifferRecords) Print() {
	for _, d := range ds {
		fmt.Printf("=> A: 0x%02x B: 0x%02x Index: 0x%x\n", d.valueA, d.valueB, d.pos)
	}
}

func Differ(a, b []uint8) DifferRecords {
	var differs []DifferRecord
	for i, v := range a {
		if v != b[i] {
			differs = append(differs, DifferRecord{
				valueA: v,
				valueB: b[i],
				pos:    i,
			})
		}
	}
	return differs
}

func ToBin(i int) string {
	s := strconv.FormatInt(int64(i), 2)
	return strings.Repeat("0", 8-len(s)) + s
}

func CompareBanks(rb []uint8, bi int, cart *cartridge.Cartridge) {
	fmt.Printf("Read bank: %d (%s)", bi, ToBin(bi))

	if Equal(rb, cart.ROMBanks[bi]) {
		fmt.Println("✔️✔️✔️✔️")
	} else {
		fe := false
		for i, b := range cart.ROMBanks {
			if Equal(b, rb) {
				fmt.Printf(" -> The bank %d (%s) matches.", i, ToBin(i))
				fe = true
				break
			}
		}
		if !fe {
			ds := Differ(cart.ROMBanks[bi], rb)
			fmt.Println("A -> Source of Truth Bank")
			fmt.Println("B -> Read bank")
			ds.Print()
		}

		fmt.Println("❌❌❌❌")
	}
	fmt.Println(strings.Repeat("*", 40))
}

func main() {
	// Read CLI flags
	mpath := flag.String("mapping", "mapping.yaml", "GameBoy pins connections to RaspberryPi pins")
	outpath := flag.String("output", "rom.gb", "Dump output file")
	flag.Parse()

	fmt.Printf("Dumping GB cartridge to %v\n", *outpath)

	// Read the RPi connections mapping to the GameBoy cartridge
	gbcon := conmap.ParseRaspberryWireMapping(*mpath)

	// Create cartridge dumper
	gbproxy := gbproxy.NewRPiGameBoyProxy(gbcon, true)
	defer gbproxy.End()
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-channel
		fmt.Println("Graceful exit dumping...", sig)
		gbproxy.End()
		os.Exit(1)
	}()

	// ofrr := cartridge.NewFileROMReader("pkmn-red.gb")
	ofrr := cartridge.NewFileROMReader("good-rom.gbc")
	cart, err := ofrr.ReadCartridge()
	if err != nil {
		fmt.Println(fmt.Errorf("file is not valid: %v", err))
		os.Exit(1)
	}
	cart.Header.PrintInfo()

	d := NewDumper(gbproxy)

	// Start reading the cartridge
	ch := d.ReadHeader()
	ch.PrintInfo()

	h := d.ReadRange(0x0, 0x150)
	if Equal(h, cart.ROMBanks[0][:0x150]) {
		fmt.Println("✔️ Raw Headers match")
	} else {
		fmt.Println("❌ Raw Headers do NOT match")
		ds := Differ(cart.ROMBanks[0][:0x150], h)
		fmt.Println("A -> Source of Truth Bank")
		fmt.Println("B -> Read bank")
		ds.Print()
	}

	err = ch.Validate()
	if err != nil {
		fmt.Println(fmt.Errorf("cartridge is not valid: %v", err))
		os.Exit(1)
	}

	nRoms := ch.GetNumROMBanks()
	banks := make([][]uint8, nRoms)
	var addrBase uint
	addrBase = 0x0
	for i := 0; i < nRoms; i++ {
		fmt.Printf("Switching to bank: 0x%02x\n", uint8(i))
		var err error
		err = d.ChangeROMBank(uint(i))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		banks[i] = d.ReadRange(addrBase, addrBase+0x4000)
		addrBase = 0x4000
		printByteArray(banks[i][:40], 4)
		CompareBanks(banks[i], i, cart)
	}

	c := cartridge.NewCartridge(ch, banks)

	err = c.Validate()
	if err != nil {
		fmt.Println(fmt.Errorf("cartridge is not valid: %v", err))
		os.Exit(1)
	}

	err = c.Save(*outpath)
	if err != nil {
		fmt.Println(fmt.Errorf("serializing dumped rom: %v", err))
		os.Exit(1)
	}

	fmt.Println("NINTENDO Logo")
	printByteArray(ch.NintendoLogo, 3)

	fmt.Println("Global Checksum")
	printByteArray(ch.GlobalChecksum, 1)
	os.Exit(0)
}
