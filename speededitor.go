package main

import (
	"log"
	"time"

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
}

type SpeedEditor struct {
	device     *hid.Device
	deviceInfo hid.DeviceInfo

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

	go func() {
		for {
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
