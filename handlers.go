package speedEditor

import (
	"github.com/JamesBalazs/speed-editor-client/hardware/keys"
	inputReport "github.com/JamesBalazs/speed-editor-client/input_report"
)

func defaultJogHandler(client SpeedEditorInt, report inputReport.JogReport) {

}

func defaultBatteryHandler(client SpeedEditorInt, report inputReport.BatteryReport) {

}

func defaultKeyPressHandler(client SpeedEditorInt, report inputReport.KeyPressReport) {
	for _, key := range report.Keys {
		if key.Led != keys.LED_NONE {
			client.SetLeds([]uint32{key.Led})
		}
		if key.JogLed != keys.LED_NONE {
			client.SetJogLeds([]uint8{key.JogLed})
		}
	}
}

func NullJogHandler(client SpeedEditorInt, report inputReport.JogReport)           {}
func NullBatteryHandler(client SpeedEditorInt, report inputReport.BatteryReport)   {}
func NullKeyPressHandler(client SpeedEditorInt, report inputReport.KeyPressReport) {}
