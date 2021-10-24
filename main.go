package main

import (
	"flag"
	"fmt"

	"github.com/Guillem96/gb-dumper/cartridge"
	"github.com/Guillem96/gb-dumper/gbproxy"
	"github.com/Guillem96/gb-dumper/io"
)

func main() {
	// Read CLI flags
	mpath := flag.String("mapping", "mapping.yaml", "GameBoy pins connections to RaspberryPi pins")
	outpath := flag.String("output", "rom.gb", "Dump output file")
	flag.Parse()

	fmt.Printf("Dumping GB cartridge to %v\n", *outpath)

	// Read the RPi connections mapping to the GameBoy cartridge
	gbcon := io.ParseWireMapping(*mpath)

	// Create cartridge dumper
	gbproxy := gbproxy.NewRPiGameBoyProxy(gbcon)
	dumper := cartridge.NewDumper(gbproxy)
	dumper.Read(0x3FFF)
}
