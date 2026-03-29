package main

import (
	"fmt"
	"log"
	"math"

	speedEditor "github.com/JamesBalazs/speed-editor-client"
	"github.com/JamesBalazs/speed-editor-client/input"
	jogModes "github.com/JamesBalazs/speed-editor-client/jog_modes"
	"github.com/JamesBalazs/speed-editor-client/keys"
	"github.com/itchyny/volume-go"
	"github.com/sstallion/go-hid"
)

var (
	keysByRow = keys.ByRow()
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

	if err := client.Authenticate(); err != nil {
		log.Fatal(err)
	}

	if err := client.SetJogMode(jogModes.ID_ABSOLUTE); err != nil {
		log.Fatal(err)
	}
	client.SetJogHandler(customJogHandler)
	client.SetKeyPressHandler(speedEditor.NullKeyPressHandler)

	client.Poll()
}

func customJogHandler(client speedEditor.SpeedEditorInt, report input.JogReport) {
	percent := (float64(report.Value) + jogModes.ABSOLUTE_MAX) / (jogModes.ABSOLUTE_MAX * 2)
	setLeds(client, percent)
	vol := int(math.Ceil(percent * 100))
	volume.SetVolume(vol)
}

func setLeds(client speedEditor.SpeedEditorInt, percent float64) {
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

	if err := client.SetLeds(leds); err != nil {
		log.Printf("error setting LEDs: %v", err)
	}
	if err := client.SetJogLeds(jogLeds); err != nil {
		log.Printf("error setting jog LEDs: %v", err)
	}
}
