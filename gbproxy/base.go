package gbproxy

type GameBoyPin interface {
	Read() bool
	High()
	Low()
	SetState(state bool)
	Input()
	Output()
}

type GameBoyProxy interface {
	Read() uint8
	SelectAddress(addr uint16)
}
