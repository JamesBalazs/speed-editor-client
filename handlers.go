package speedEditor

import (
	"fmt"

	"github.com/JamesBalazs/speed-editor-client/input"
)

func defaultJogHandler(client SpeedEditorInt, report input.JogReport) {
	fmt.Printf("Jog mode: %s, value: %d\n", report.Mode.Name, report.Value)
}

func defaultBatteryHandler(client SpeedEditorInt, report input.BatteryReport) {
	fmt.Printf("Battery charging: %v, level: %f\n", report.Charging, report.Battery)
}

func defaultKeyPressHandler(client SpeedEditorInt, report input.KeyPressReport) {
	for _, key := range report.Keys {
		fmt.Printf("Keys pressed: %s", key.Name)
	}
}

func NullJogHandler(client SpeedEditorInt, report input.JogReport)           {}
func NullBatteryHandler(client SpeedEditorInt, report input.BatteryReport)   {}
func NullKeyPressHandler(client SpeedEditorInt, report input.KeyPressReport) {}
