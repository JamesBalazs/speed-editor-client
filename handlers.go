package speedEditor

import (
	"github.com/JamesBalazs/speed-editor-client/hardware/keys"
	inputReport "github.com/JamesBalazs/speed-editor-client/input_report"
)

func defaultJogHandler(se *SpeedEditor, report inputReport.JogReport) {

}

func defaultBatteryHandler(se *SpeedEditor, report inputReport.BatteryReport) {

}

func defaultKeyPressHandler(se *SpeedEditor, report inputReport.KeyPressReport) {
	for _, key := range report.Keys {
		if key.Led != keys.LED_NONE {
			se.SetLeds([]uint32{key.Led})
		}
		if key.JogLed != keys.LED_NONE {
			se.SetJogLeds([]uint8{key.JogLed})
		}
	}
}
