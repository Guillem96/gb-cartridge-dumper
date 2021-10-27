package gbproxy

// GameBoyPin interface defines the needed methods to manage the GameBoy connections
type GameBoyPin interface {
	// Read returns the pin state, if returns true it means that the pin state is High
	Read()

	// High sets the pin value to 1
	High()

	// Low sets the pin status to 0
	Low()

	// SetState sets the given state to the pin, if true is provided as a parameter the pin state
	// will be set to High
	SetState(state bool)

	// Input sets the pin in input mode (pins controlled from the host)
	Input()

	// Output sets the pin in output mode (pin status is populated by the GameBoy)
	Output()
}

// Interface to read and write data to GameBoy pins from your hardware
// The cartridge.Dumper depends on this interface. In this project I am using a RaspberryPi and
// I am implementing this interface in the gbproxy/gbrpi.go so it works with it.
// Ideally, if you are working with any other type of hardware (arduino for instance) you should
// only implement this interface, so you have an object able to interact with GameBoy, and provide
// this new object to the cartrige.Dumper
type GameBoyProxy interface {

	// Read returns the byte read at the address specified with the SelectAddress method
	Read() uint8

	// SelectAddress sets the pins status so they point to the provided memory address
	SelectAddress(addr uint16)
}
