package main

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"time"

	speedEditor "github.com/JamesBalazs/speed-editor-client"
	"github.com/JamesBalazs/speed-editor-client/hardware/keys"
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

	index := 0
	leds := slices.Collect(maps.Values(keys.Leds))
	for {
		time.Sleep(75 * time.Millisecond)
		client.SetLeds(leds[0:index])

		if index > len(leds) {
			index = 0
		} else {
			index += 1
		}
	}
}
