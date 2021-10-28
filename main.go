package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Guillem96/gb-dumper/cartridge"
	"github.com/Guillem96/gb-dumper/gbproxy"
	"github.com/Guillem96/gb-dumper/io"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	// Read CLI flags
	mpath := flag.String("mapping", "mapping.yaml", "GameBoy pins connections to RaspberryPi pins")
	outpath := flag.String("output", "rom.gb", "Dump output file")
	flag.Parse()

	fmt.Printf("Dumping GB cartridge to %v\n", *outpath)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	// Read the RPi connections mapping to the GameBoy cartridge
	gbcon := io.ParseWireMapping(*mpath)

	// Create cartridge dumper
	gbproxy := gbproxy.NewRPiGameBoyProxy(gbcon)
	d := cartridge.NewDumper(gbproxy)

	// Start reading the cartridge
	ch, err := d.ReadHeader()
	if err != nil {
		fmt.Println(fmt.Errorf("reading header: %v\n", err))
		os.Exit(1)
	}

	err = ch.ValidateHeader()
	if err != nil {
		fmt.Println(fmt.Errorf("cartridge is not valid: %v", err))
		os.Exit(1)
	}

	fmt.Println(ch.NintendoLogo)
	os.Exit(0)
}
