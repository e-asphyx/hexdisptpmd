package main

import (
	"github.com/e-asphyx/tpm2net"
	"log"
)

type TPMHandler struct {
	VirtualWidth  int
	VirtualHeight int

	Oversample int

	Layout Layout
	Output OutputDevice

	frameBuffer  []uint8
	outputBuffer []Color
	cellMap      []map[int]float64
}

func (hnd *TPMHandler) HandlePacket(pkt *tpm2net.Packet) {
	numLeds := hnd.Layout.NumCells()
	bufSz := hnd.VirtualWidth * hnd.VirtualHeight * 3

	if len(hnd.frameBuffer) < bufSz {
		hnd.frameBuffer = make([]uint8, bufSz)
	}
	copy(hnd.frameBuffer, pkt.Data)

	if len(hnd.outputBuffer) < numLeds {
		hnd.outputBuffer = make([]Color, numLeds)
	}

	if len(hnd.cellMap) < numLeds {
		hnd.initMap()
	}

	for i := range hnd.outputBuffer {
		acc := [3]float64{}

		for offs, coeff := range hnd.cellMap[i] {
			acc[0] += float64(hnd.frameBuffer[offs*3]) * coeff
			acc[1] += float64(hnd.frameBuffer[offs*3+1]) * coeff
			acc[2] += float64(hnd.frameBuffer[offs*3+2]) * coeff
		}

		for c := range acc {
			if acc[c] > 255.0 {
				acc[c] = 255.0
			}
		}

		hnd.outputBuffer[i] = RGB(uint8(acc[0]+0.5), uint8(acc[1]+0.5), uint8(acc[2]+0.5))
	}

	if _, err := hnd.Output.Write(hnd.outputBuffer); err != nil {
		log.Println(err)
	}
}

func (hnd *TPMHandler) initMap() {
	hnd.cellMap = make([]map[int]float64, hnd.Layout.NumCells())

	layoutWidth := hnd.Layout.Width()
	layoutHeight := hnd.Layout.Height()

	virtualWidth := hnd.VirtualWidth * hnd.Oversample
	virtualHeight := hnd.VirtualHeight * hnd.Oversample

	xstep := layoutWidth / float64(virtualWidth)
	ystep := layoutHeight / float64(virtualHeight)

	for i := range hnd.cellMap {
		poly := hnd.Layout.Cell(i)
		tl, br := poly.BoundingBox()

		x0 := int(tl.X / layoutWidth * float64(virtualWidth))
		x1 := int(br.X / layoutWidth * float64(virtualWidth))

		y0 := int(tl.Y / layoutHeight * float64(virtualHeight))
		y1 := int(br.Y / layoutHeight * float64(virtualHeight))

		cnt := make(map[int]int)
		sum := 0

		for y := y0; y <= y1; y++ {
			for x := x0; x <= x1; x++ {
				center := Point{float64(x)*xstep + xstep/2.0, float64(y)*ystep + ystep/2.0}

				if poly.PointInside(center) &&
					x >= 0 && y >= 0 &&
					x < virtualWidth && y < virtualHeight {

					xx := x / hnd.Oversample
					yy := y / hnd.Oversample
					offs := (yy*hnd.VirtualWidth + xx)

					cnt[offs]++
					sum++
				}
			}
		}

		hnd.cellMap[i] = make(map[int]float64)
		for offs, c := range cnt {
			hnd.cellMap[i][offs] = float64(c) / float64(sum)
		}
	}
}
