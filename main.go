package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Guillem96/gameboy-tools/conmap"
	"github.com/Guillem96/gameboy-tools/gbproxy"
)

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
	d := NewDumper(gbproxy)

	// Start reading the cartridge
	ch := d.ReadHeader()

	err := ch.ValidateHeader()
	if err != nil {
		fmt.Println(fmt.Errorf("cartridge is not valid: %v", err))
		os.Exit(1)
	}

	fmt.Println(ch.NintendoLogo)
	os.Exit(0)
}
