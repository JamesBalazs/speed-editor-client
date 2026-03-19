package main

import (
	"fmt"
	"log"

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

	for {
		data, len := speedEditor.Read()

		fmt.Printf("len: %d data: %v\n", len, data)
	}
}
