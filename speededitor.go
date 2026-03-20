package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	inputReport "github.com/JamesBalazs/speed-editor-client/input_report"
	"github.com/sstallion/go-hid"
)

const VID = 0x1edb
const PID = 0xda0e

// NewSpeedEditor connects to a Speed Editor via the HID library
// and returns a SpeedEditorInt to interact with the device.
//
// It is recommended to manually initialise the HID library before
// creating the Speed Editor client, with `hid.Init()`.
//
// Ensure to use `defer hid.Exit()` to avoid memory leaks.
func NewSpeedEditor() SpeedEditorInt {
	device, err := hid.OpenFirst(VID, PID)
	if err != nil {
		log.Fatal(err)
	}

	speedEditor := &SpeedEditor{
		device:      device,
		AuthHandler: AuthHandler{device},
	}
	speedEditor.initialize()

	return speedEditor
}

type SpeedEditorInt interface {
	// Authenticate does the initial handshake with the Speed Editor,
	// and re-auths periodically in the background when requested by the device.
	Authenticate()

	// GetDeviceInfo returns the serial number, manufacturer string etc published
	// by the device via HID. This info is cached on init, so we don't have to
	// request it on every call.
	GetDeviceInfo() hid.DeviceInfo

	// Read pulls a single input report from the Speed Editor, and returns the
	// data (list of keys currently held down, jog wheel position etc.) and
	// the data length.
	//
	// The first byte indicates which report type was received.
	Read() ([]byte, int)

	Poll()

	SetLeds(leds []uint32)
}

type SpeedEditor struct {
	device     *hid.Device
	deviceInfo hid.DeviceInfo
	activeLeds []uint32

	AuthHandler AuthHandlerInt
}

// initialize grabs the device's serial number, manufacturer string etc via HID.
// The handshake is not required before this step can take place.
func (se *SpeedEditor) initialize() {
	deviceInfo, err := se.device.GetDeviceInfo()
	if err != nil {
		log.Fatal(err)
	}

	se.deviceInfo = *deviceInfo
}

func (se SpeedEditor) Authenticate() {
	// Getting the initial reAuthSeconds synchronously.
	// Do not read or update this outside of the goroutine to avoid a data race.
	reAuthSeconds := se.AuthHandler.Authenticate()

	fmt.Printf("Initial handshake\n")

	go func() {
		for {
			fmt.Printf("Sleeping %s\n", reAuthSeconds)

			time.Sleep(reAuthSeconds)

			reAuthSeconds = se.AuthHandler.Authenticate()
		}
	}()
}

func (se SpeedEditor) GetDeviceInfo() hid.DeviceInfo {
	return se.deviceInfo
}

func (se SpeedEditor) Read() ([]byte, int) {
	data := make([]byte, 16)
	len, err := se.device.Read(data)

	if err != nil {
		log.Fatal(err)
	}

	return data, len
}

func (se SpeedEditor) Poll() {
	for {
		data, len := se.Read()
		report := inputReport.NewInputReport(data, len)
		report.Handle()
	}
}

const LedReportId = 2

func (se SpeedEditor) SetLeds(leds []uint32) {
	payload := make([]byte, 5)
	payload[0] = LedReportId

	var bitMask uint32
	for _, led := range leds {
		bitMask |= led
	}
	binary.LittleEndian.PutUint32(payload[1:], bitMask)

	se.device.Write(payload)
}
