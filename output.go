package main

type Color uint32

type OutputDevice interface {
	Write(leds []Color) (n int, err error)
}

func RGB(r, g, b uint8) Color {
	return Color((uint32(r) << 16) | (uint32(g) << 8) | uint32(b))
}

func (c Color) R() uint8 {
	return uint8((uint32(c) >> 16) & 0xff)
}

func (c Color) G() uint8 {
	return uint8((uint32(c) >> 8) & 0xff)
}

func (c Color) B() uint8 {
	return uint8(uint32(c) & 0xff)
}
