package main

import (
	"github.com/e-asphyx/rpiws"
	"log"
)

type WSOutput struct {
	driver *rpiws.Driver
	chnum  int
	ch     chan []Color
	notify chan int
}

func NewWSOutput(driver *rpiws.Driver, channel int) *WSOutput {
	out := WSOutput{
		driver: driver,
		chnum:  channel,
		ch:     make(chan []Color, 1),
		notify: make(chan int),
	}
	go out.serve()
	return &out
}

func (out *WSOutput) serve() {
	for {
		buf, ok := <-out.ch
		if !ok {
			break
		}

		leds := out.driver.Channel[out.chnum].Leds()
		for i, led := range buf {
			if i < len(leds) {
				leds[i] = rpiws.RGB(led.R(), led.G(), led.B())
			}
		}

		if err := out.driver.Render(); err != nil {
			log.Println(err)
			break
		}
	}
	out.notify <- 1
}

func (out *WSOutput) Write(leds []Color) (n int, err error) {
	if len(out.ch) < cap(out.ch) {
		out.ch <- leds
		return len(leds), nil
	} else {
		return 0, nil
	}
}

func (out *WSOutput) Stop() {
	close(out.ch)
	<-out.notify
}
