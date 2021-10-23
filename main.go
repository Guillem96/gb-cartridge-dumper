package main

import (
	"flag"
	"fmt"

	"github.com/Guillem96/gb-dumper/cartridge"
	"github.com/Guillem96/gb-dumper/io"
)

func main() {
	// Read CLI flags
	mappingPath := flag.String("mapping", "mapping.yaml", "GameBoy pins connections to RaspberryPi pins")
	outputPath := flag.String("output", "rom.gb", "Dump output file")
	flag.Parse()

	fmt.Printf("Dumping GB cartridge to %v\n", *outputPath)

	// Read the RPi connections mapping to the GameBoy cartridge
	gbConnections := io.ParseWireMapping(*mappingPath)

	// Create cartridge dumper
	dumper := cartridge.NewDumper(gbConnections)
	dumper.ReadAddress(0x3FFF)
}
