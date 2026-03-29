package speedEditor

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/JamesBalazs/speed-editor-client/input"
	"github.com/sstallion/go-hid"
)

const (
	VID = 0x1edb
	PID = 0xda0e

	LedReportId     = 2
	JogModeReportId = 3
	JogLedReportId  = 4
)

// deviceInterface defines the HID device operations for testability
type deviceInterface interface {
	Close() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	GetDeviceInfo() (*hid.DeviceInfo, error)
	GetFeatureReport(buf []byte) (int, error)
	SendFeatureReport(buf []byte) (int, error)
}

// hidDeviceWrapper wraps the hid.Device to implement deviceInterface
type hidDeviceWrapper struct {
	device *hid.Device
}

func (w *hidDeviceWrapper) Close() error {
	return w.device.Close()
}

func (w *hidDeviceWrapper) Read(buf []byte) (int, error) {
	return w.device.Read(buf)
}

func (w *hidDeviceWrapper) Write(buf []byte) (int, error) {
	return w.device.Write(buf)
}

func (w *hidDeviceWrapper) GetDeviceInfo() (*hid.DeviceInfo, error) {
	return w.device.GetDeviceInfo()
}

func (w *hidDeviceWrapper) GetFeatureReport(buf []byte) (int, error) {
	return w.device.GetFeatureReport(buf)
}

func (w *hidDeviceWrapper) SendFeatureReport(buf []byte) (int, error) {
	return w.device.SendFeatureReport(buf)
}

// NewClient connects to a Speed Editor via the HID library
// and returns a SpeedEditorInt to interact with the device.
//
// It is recommended to manually initialise the HID library before
// creating the Speed Editor client, with `hid.Init()`.
//
// Ensure to use `defer hid.Exit()` to avoid memory leaks.
func NewClient() (SpeedEditorInt, error) {
	device, err := hid.OpenFirst(VID, PID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	wrapper := &hidDeviceWrapper{device: device}

	speedEditor := &SpeedEditor{
		device:      wrapper,
		AuthHandler: AuthHandler{device: wrapper},
	}

	if err = speedEditor.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	return speedEditor, nil
}

type SpeedEditorInt interface {
	// Authenticate does the initial handshake with the Speed Editor,
	// and re-auths periodically in the background when requested by the device.
	Authenticate() error

	// GetDeviceInfo returns the serial number, manufacturer string etc published
	// by the device via HID. This info is cached on init, so we don't have to
	// request it on every call.
	GetDeviceInfo() hid.DeviceInfo

	// Read pulls a single input report from the Speed Editor, and returns the
	// data (list of keys currently held down, jog wheel position etc.) and
	// the data length.
	//
	// The first byte indicates which report type was received.
	Read() ([]byte, int, error)

	// Poll starts a Read loop, parses each input report and calls Handle on each
	// via the Report interface.
	Poll()

	// SetLeds accepts the bitmask for a list of LEDs, and binary ORs the bitmask
	// to enable all LEDs in the mask.
	//
	// SetLeds does not keep any state, so it will reset any previously enabled
	// LEDs if they aren't included in the next call.
	SetLeds(leds []uint32) error

	// SetJogMode switches between the 4 jog modes:
	// RELATIVE - Relative position
	// ABSOLUTE - Absolute position from -4096 to 4096
	// RELATIVE2 - Relative position, I think this is used to enable a faster scroll mode when the SCRL button is pressed twice in Resolve: https://www.reddit.com/r/blackmagicdesign/comments/1dv56d4/speed_editor_firmware_update_dial_speed_change/
	// ABSOLUTE_0 - Absolute position from -4096 to 4096 with deadzone around 0
	SetJogMode(mode uint8) error

	// SetJogLeds accepts the bitmask for a list of LEDs, and binary ORs the bitmask
	// to enable all LEDs in the mask. Jog LEDs are on a separate system, and overlap
	// with some of the normal LED IDs so we need a separate message to enable them.
	//
	// SetJogLeds does not keep any state, so it will reset any previously enabled
	// LEDs if they aren't included in the next call.
	SetJogLeds(leds []uint8) error

	// SetJogHandler allows replacing the handler function that will be called on Poll()
	// when a JogReport is received.
	SetJogHandler(handler func(SpeedEditorInt, input.JogReport))
	// SetBatteryHandler allows replacing the handler function that will be called on Poll()
	// when a BatteryReport is received.
	SetBatteryHandler(handler func(SpeedEditorInt, input.BatteryReport))
	// SetKeyPressHandler allows replacing the handler function that will be called on Poll()
	// when a KeyPressReport is received.
	SetKeyPressHandler(handler func(SpeedEditorInt, input.KeyPressReport))
}

type SpeedEditor struct {
	device     deviceInterface
	deviceInfo hid.DeviceInfo
	activeLeds []uint32

	AuthHandler AuthHandlerInt

	jogHandler      func(SpeedEditorInt, input.JogReport)
	batteryHandler  func(SpeedEditorInt, input.BatteryReport)
	keyPressHandler func(SpeedEditorInt, input.KeyPressReport)
}

// initialize grabs the device's serial number, manufacturer string etc via HID.
// The handshake is not required before this step can take place.
func (se *SpeedEditor) initialize() error {
	deviceInfo, err := se.device.GetDeviceInfo()
	if err != nil {
		return fmt.Errorf("failed to get device info: %w", err)
	}

	se.deviceInfo = *deviceInfo

	se.SetJogHandler(defaultJogHandler)
	se.SetBatteryHandler(defaultBatteryHandler)
	se.SetKeyPressHandler(defaultKeyPressHandler)

	return nil
}

func (se SpeedEditor) Authenticate() error {
	// Getting the initial reAuthSeconds synchronously.
	// Do not read or update this outside of the goroutine to avoid a data race.
	reAuthSeconds, err := se.AuthHandler.Authenticate()
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	go func() {
		for {
			time.Sleep(reAuthSeconds)

			reAuthSeconds, err = se.AuthHandler.Authenticate()
			if err != nil {
				fmt.Printf("failed to re-authenticate: %v\n", err)
			}
		}
	}()

	return nil
}

func (se SpeedEditor) GetDeviceInfo() hid.DeviceInfo {
	return se.deviceInfo
}

func (se SpeedEditor) Read() ([]byte, int, error) {
	data := make([]byte, 9)
	n, err := se.device.Read(data)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to read from device: %w", err)
	}

	return data, n, nil
}

func (se SpeedEditor) Poll() {
	for {
		data, _, err := se.Read()
		if err != nil {
			fmt.Printf("error reading from device: %v\n", err)
			continue
		}
		report, parseErr := input.ReportBytes(data).ToReport()
		if parseErr != nil {
			fmt.Printf("error parsing report: %v\n", parseErr)
			continue
		}
		se.HandleReport(report)
	}
}

func (se SpeedEditor) HandleReport(genericReport any) {
	switch report := genericReport.(type) {
	case input.JogReport:
		se.HandleJog(report)
	case input.BatteryReport:
		se.HandleBattery(report)
	case input.KeyPressReport:
		se.HandleKeyPress(report)
	} // TODO handle unknown reports (log error and continue)
}

func (se SpeedEditor) SetLeds(leds []uint32) error {
	payload := make([]byte, 5)
	payload[0] = LedReportId

	var bitField uint32
	for _, bitMask := range leds {
		bitField |= bitMask
	}
	binary.LittleEndian.PutUint32(payload[1:], bitField)

	_, err := se.device.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to set LEDs: %w", err)
	}

	return nil
}

func (se SpeedEditor) SetJogMode(mode uint8) error {
	payload := make([]byte, 7)
	payload[0] = JogModeReportId
	payload[1] = mode
	// bytes 3-6 are zero
	payload[6] = 255 // byte 7 has unknown purpose

	_, err := se.device.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to set jog mode: %w", err)
	}

	return nil
}

func (se SpeedEditor) SetJogLeds(leds []uint8) error {
	payload := make([]byte, 2)
	payload[0] = JogLedReportId

	var bitField uint8
	for _, bitMask := range leds {
		bitField |= bitMask
	}
	payload[1] = bitField

	_, err := se.device.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to set jog LEDs: %w", err)
	}

	return nil
}

func (se *SpeedEditor) SetJogHandler(handler func(SpeedEditorInt, input.JogReport)) {
	se.jogHandler = handler
}

func (se *SpeedEditor) SetBatteryHandler(handler func(SpeedEditorInt, input.BatteryReport)) {
	se.batteryHandler = handler
}

func (se *SpeedEditor) SetKeyPressHandler(handler func(SpeedEditorInt, input.KeyPressReport)) {
	se.keyPressHandler = handler
}

func (se SpeedEditor) HandleJog(report input.JogReport) {
	se.jogHandler(&se, report)
}

func (se SpeedEditor) HandleBattery(report input.BatteryReport) {
	se.batteryHandler(&se, report)
}

func (se SpeedEditor) HandleKeyPress(report input.KeyPressReport) {
	se.keyPressHandler(&se, report)
}
