package input

import (
	"fmt"

	"github.com/JamesBalazs/speed-editor-client/keys"
)

func init() {
	// Cache constant-time key index so we don't copy-on-read
	// every time a handler wants to look up a key.
	keysById = keys.ById()
}

const (
	ReportJog      = 3
	ReportKeyPress = 4
	ReportBattery  = 7
)

var (
	keysById = map[uint16]keys.Key{}
)

type ReportBytes []byte

func (byt ReportBytes) ToReport() (any, error) {
	id := byt[0]
	payload := byt[1:]

	switch id {
	case ReportJog:
		return NewJogReport(id, payload), nil
	case ReportKeyPress:
		return NewKeyPressReport(id, payload)
	case ReportBattery:
		return NewBatteryReport(id, payload)
	default:
		return nil, fmt.Errorf("unknown report id: %v", id)
	}
}

type ReportInt interface {
	Handle()
}
