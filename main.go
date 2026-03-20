package main

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"time"

	"github.com/JamesBalazs/speed-editor-client/hardware/keys"
	"github.com/sstallion/go-hid"
)

func main() {
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}
	defer hid.Exit()

	speedEditor := NewSpeedEditor()

	deviceInfo := speedEditor.GetDeviceInfo()

	fmt.Printf("Manufacturer: %s\nProduct: %s\nSerial: %s\n", deviceInfo.MfrStr, deviceInfo.ProductStr, deviceInfo.SerialNbr)

	speedEditor.Authenticate()
	// speedEditor.Poll()
	index := 0
	leds := slices.Collect(maps.Values(keys.Leds))
	for {
		time.Sleep(100 * time.Millisecond)
		speedEditor.SetLeds(leds[0:index])

		if index > len(leds) {
			index = 0
		} else {
			index += 1
		}
	}
}
