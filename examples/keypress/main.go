package main

import (
	"fmt"
	"log"

	speedEditor "github.com/JamesBalazs/speed-editor-client"
	"github.com/JamesBalazs/speed-editor-client/input"
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

	client.SetKeyPressHandler(customKeyPressHandler)

	client.Poll()
}

func customKeyPressHandler(client speedEditor.SpeedEditorInt, report input.KeyPressReport) {
	for _, key := range report.Keys {
		if key.Led != keys.LED_NONE {
			client.SetLeds([]uint32{key.Led})
		}
		if key.JogLed != keys.LED_NONE {
			client.SetJogLeds([]uint8{key.JogLed})
		}
	}
}
