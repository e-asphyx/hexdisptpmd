package main

import (
	"github.com/e-asphyx/tpm2net"
	"log"
)

type TPMHandler struct {
	VirtualWidth  int
	VirtualHeight int

	Layout Layout
	Output OutputDevice

	frameBuffer  []uint8
	outputBuffer []Color
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

	layoutWidth := hnd.Layout.Width()
	layoutHeight := hnd.Layout.Height()

	xstep := layoutWidth / float64(hnd.VirtualWidth)
	ystep := layoutHeight / float64(hnd.VirtualHeight)

	for i := range hnd.outputBuffer {
		poly := hnd.Layout.Cell(i)
		tl, br := poly.BoundingBox()

		x0 := int(tl.X / layoutWidth * float64(hnd.VirtualWidth))
		x1 := int(br.X / layoutWidth * float64(hnd.VirtualWidth))

		y0 := int(tl.Y / layoutHeight * float64(hnd.VirtualHeight))
		y1 := int(br.Y / layoutHeight * float64(hnd.VirtualHeight))

		acc := [3]int32{}
		var cnt int32

		for y := y0; y <= y1; y++ {
			for x := x0; x <= x1; x++ {
				center := Point{float64(x)*xstep + xstep/2.0, float64(y)*ystep + ystep/2.0}

				if poly.PointInside(center) &&
					x >= 0 && y >= 0 &&
					x < hnd.VirtualWidth && y < hnd.VirtualHeight {

					offs := (y*hnd.VirtualWidth + x) * 3
					acc[0] += int32(hnd.frameBuffer[offs])
					acc[1] += int32(hnd.frameBuffer[offs+1])
					acc[2] += int32(hnd.frameBuffer[offs+2])
					cnt++
				}
			}
		}

		if cnt != 0 {
			for c := range acc {
				acc[c] = (acc[c] + (cnt / 2)) / cnt
				if acc[c] > 255 {
					acc[c] = 255
				}
			}
		} else {
			log.Println("Misplaced cell")
		}

		hnd.outputBuffer[i] = RGB(uint8(acc[0]), uint8(acc[1]), uint8(acc[2]))
	}

	if _, err := hnd.Output.Write(hnd.outputBuffer); err != nil {
		log.Println(err)
	}
}
