package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Guillem96/gameboy-tools/conmap"
	"github.com/Guillem96/gameboy-tools/gbproxy"
)

func main() {
	// Read CLI flags
	mpath := flag.String("mapping", "mapping.yaml", "GameBoy pins connections to RaspberryPi pins")
	outpath := flag.String("output", "rom.gb", "Dump output file")
	skipChecksStr := flag.String("skip-checks", "no", "Wether to skip the checksum checks or not")
	flag.Parse()

	*skipChecksStr = strings.ToLower(*skipChecksStr)
	sc := *skipChecksStr == "y" || *skipChecksStr == "yes" || *skipChecksStr == "true" || *skipChecksStr == "1"

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

	d := NewDumper(gbproxy)

	// Dump the header and validate the checkpoint
	ch := d.ReadHeader()
	ch.PrintInfo()
	err := ch.Validate()
	if err != nil && sc {
		fmt.Println(fmt.Errorf("cartridge is not valid: %+v", err))
		os.Exit(1)
	} else if err != nil && !sc {
		fmt.Printf("Warning: cartridge is not valid: %+v\n", err)
	} else {
		fmt.Println("✔️ Header is valid!")
	}

	// Dump the whole cartridge
	c, err := d.ReadCartridge()
	if err != nil {
		fmt.Println(fmt.Errorf("error dumping the cartridge: %v", err))
		os.Exit(1)
	}

	// Validate the global checksum
	err = c.Validate()
	if err != nil && sc {
		fmt.Println(fmt.Errorf("cartridge is not valid: %v", err))
		os.Exit(1)
	} else if err != nil && !sc {
		fmt.Printf("Warning: cartridge is not valid: %+v\n", err)
	} else {
		fmt.Println("✔️ Cartridge global checksum is valid!")
	}

	err = c.Save(*outpath)
	if err != nil {
		fmt.Println(fmt.Errorf("serializing dumped rom: %v", err))
		os.Exit(1)
	}

	os.Exit(0)
}
