package main

import (
	"fmt"
	"log"

	"github.com/sstallion/go-hid"

	inputReport "github.com/JamesBalazs/speed-editor-client/input_report"
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
		report := inputReport.NewInputReport(data, len)
		report.Handle()
	}
}
