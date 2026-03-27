package main

import (
	"fmt"
	"log"
	"time"

	speedEditor "github.com/JamesBalazs/speed-editor-client"
	"github.com/JamesBalazs/speed-editor-client/keys"
	"github.com/sstallion/go-hid"
)

func main() {
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}
	defer hid.Exit()

	client := speedEditor.NewClient()

	deviceInfo := client.GetDeviceInfo()

	fmt.Printf("Manufacturer: %s\nProduct: %s\nSerial: %s\n", deviceInfo.MfrStr, deviceInfo.ProductStr, deviceInfo.SerialNbr)

	client.Authenticate()

	keysByCol := keys.ByCol()
	keysByRow := keys.ByRow()
	for {
		for x := range 10 {
			column := keysByCol[float32(x)]
			leds := []uint32{}
			jogLeds := []uint8{}
			for _, key := range column {
				if x >= 7 {
					jogLeds = append(jogLeds, key.JogLed)
				} else {
					leds = append(leds, key.Led)
				}
			}
			client.SetLeds(leds)
			client.SetJogLeds(jogLeds)
			time.Sleep(75 * time.Millisecond)
		}

		for y := range 6 {
			row := keysByRow[int(y)]
			leds := []uint32{}
			jogLeds := []uint8{}
			for x, key := range row {
				if x >= 7 {
					jogLeds = append(jogLeds, key.JogLed)
				} else {
					leds = append(leds, key.Led)
				}
			}
			client.SetLeds(leds)
			client.SetJogLeds(jogLeds)
			time.Sleep(75 * time.Millisecond)
		}
	}
}
