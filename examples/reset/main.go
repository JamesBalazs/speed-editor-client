package main

import (
	"fmt"
	"log"

	speedEditor "github.com/JamesBalazs/speed-editor-client"
	"github.com/sstallion/go-hid"
)

func main() {
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}
	defer hid.Exit()

	client, err := speedEditor.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	deviceInfo := client.GetDeviceInfo()

	fmt.Printf("Manufacturer: %s\nProduct: %s\nSerial: %s\n", deviceInfo.MfrStr, deviceInfo.ProductStr, deviceInfo.SerialNbr)

	client.Authenticate()

	if err := client.SetLeds([]uint32{}); err != nil {
		log.Fatal(err)
	}
	if err := client.SetJogLeds([]uint8{}); err != nil {
		log.Fatal(err)
	}
}
