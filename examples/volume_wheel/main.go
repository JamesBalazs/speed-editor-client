package main

import (
	"fmt"
	"log"
	"math"

	speedEditor "github.com/JamesBalazs/speed-editor-client"
	"github.com/JamesBalazs/speed-editor-client/hardware"
	"github.com/JamesBalazs/speed-editor-client/hardware/keys"
	inputReport "github.com/JamesBalazs/speed-editor-client/input_report"
	"github.com/sstallion/go-hid"
)

var (
	percent float64

	keysByRow = keys.ByRow()
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

	client.SetJogMode(hardware.JOGMODE_ABSOLUTE)
	client.SetJogHandler(customJogHandler)
	client.SetKeyPressHandler(speedEditor.NullKeyPressHandler)

	client.Poll()
}

func customJogHandler(client speedEditor.SpeedEditorInt, report inputReport.JogReport) {
	percent = (float64(report.Value) + 4096) / 8192
	setLeds(client)
}

func setLeds(client speedEditor.SpeedEditorInt) {
	rows := int(math.Ceil(percent * 6))
	leds := []uint32{}
	jogLeds := []uint8{}

	for y := range rows {
		row := keysByRow[6-y]
		for x, key := range row {
			if x >= 7 {
				jogLeds = append(jogLeds, key.JogLed)
			} else {
				leds = append(leds, key.Led)
			}
		}
	}

	client.SetLeds(leds)
	client.SetJogLeds(jogLeds)
}
