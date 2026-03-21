package main

import (
	"fmt"
	"log"
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

	keysByPos := keys.ByPos()
	for {
		for x := 0; x < 10; x++ {
			keysOnColumn := keysByPos[float32(x)]
			leds := []uint32{}
			for _, key := range keysOnColumn {
				leds = append(leds, key.Led)
			}
			client.SetLeds(leds)
			time.Sleep(75 * time.Millisecond)
		}
	}
}
