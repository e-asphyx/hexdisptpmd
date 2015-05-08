package main

import (
	"github.com/e-asphyx/rpiws"
	"github.com/e-asphyx/tpm2net"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	ListenAddr    = ":65506"
	VirtualWidth  = 44
	VirtualHeight = 44
	Oversample    = 8
)

func main() {
	// Layout manager
	layout := HexLayout{
		EdgeLen:       6,
		InvertOddRows: true,
	}

	// WS281x driver
	wsDriver := rpiws.Driver{
		Freq:   rpiws.WS2811_TARGET_FREQ,
		Dmanum: 5,
		Channel: [rpiws.RPI_PWM_CHANNELS]rpiws.Channel{
			rpiws.Channel{
				Gpionum:    18,
				Count:      int32(layout.NumCells()),
				Brightness: 255,
			},
		},
	}

	err := wsDriver.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer wsDriver.Fini()

	output := NewWSOutput(&wsDriver, 0)

	handler := TPMHandler{
		VirtualWidth:  VirtualWidth,
		VirtualHeight: VirtualHeight,

		Oversample: Oversample,

		Layout: &layout,
		Output: output,
	}

	// TPM2Net server
	server := tpm2net.Server{
		MaxPacketNum:  0,
		MaxPacketSize: VirtualWidth * VirtualHeight * 3,
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

	conn, err := net.ListenPacket("udp", ListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	/*---------------------------------------------------------------*/

	log.Println("Ready")

	notify := make(chan int)
	go func() {
		server.Serve(conn, &handler)
		notify <- 1
	}()

	<-sigch
	conn.Close()
	<-notify
	output.Stop()

	log.Println("Bye")
}
