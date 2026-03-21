package speedEditor

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

// NewClient connects to a Speed Editor via the HID library
// and returns a SpeedEditorInt to interact with the device.
//
// It is recommended to manually initialise the HID library before
// creating the Speed Editor client, with `hid.Init()`.
//
// Ensure to use `defer hid.Exit()` to avoid memory leaks.
func NewClient() SpeedEditorInt {
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

	// Poll starts a Read loop, parses each input report and calls Handle on each
	// via the Report interface.
	Poll()

	// SetLeds accepts the bitmask for a list of LEDs, and binary ORs the bitmask
	// to enable all LEDs in the mask.
	//
	// SetLeds does not keep any state, so it will reset any previously enabled
	// LEDs if they aren't included in the next call.
	SetLeds(leds []uint32)

	// SetJogMode switches between the 4 jog modes:
	// RELATIVE - Relative position
	// ABSOLUTE - Absolute position from -4096 to 4096
	// RELATIVE2 - Relative position, I think this is used to enable a faster scroll mode when the SCRL button is pressed twice in Resolve: https://www.reddit.com/r/blackmagicdesign/comments/1dv56d4/speed_editor_firmware_update_dial_speed_change/
	// ABSOLUTE_0 - Absolute position from -4096 to 4096 with deadzone around 0
	SetJogMode(mode uint8)

	// SetJogLeds accepts the bitmask for a list of LEDs, and binary ORs the bitmask
	// to enable all LEDs in the mask. Jog LEDs are on a separate system, and overlap
	// with some of the normal LED IDs so we need a separate message to enable them.
	//
	// SetJogLeds does not keep any state, so it will reset any previously enabled
	// LEDs if they aren't included in the next call.
	SetJogLeds(leds []uint8)
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

	var bitField uint32
	for _, bitMask := range leds {
		bitField |= bitMask
	}
	binary.LittleEndian.PutUint32(payload[1:], bitField)

	se.device.Write(payload)
}

const JogModeReportId = 3

func (se SpeedEditor) SetJogMode(mode uint8) {
	payload := make([]byte, 7)
	payload[0] = JogModeReportId
	payload[1] = mode
	// bytes 3-6 are zero
	payload[6] = 255 // byte 7 has unknown purpose

	se.device.Write(payload)
}

const JogLedReportId = 4

func (se SpeedEditor) SetJogLeds(leds []uint8) {
	payload := make([]byte, 2)
	payload[0] = JogLedReportId

	var bitField uint8
	for _, bitMask := range leds {
		bitField |= bitMask
	}
	payload[1] = bitField

	se.device.Write(payload)
}
