package inputReport

import "fmt"

const (
	ReportJog      = 3
	ReportKeyPress = 4
	ReportBattery  = 7
)

func NewInputReport(raw []byte, len int) ReportInt {
	id := raw[0]
	payload := raw[1:len]
	len = len - 1

	switch id {
	case ReportJog:
		return NewJogReport(id, payload, len)
	case ReportKeyPress:
		return NewKeyPressReport(id, payload, len)
	case ReportBattery:
		return NewBatteryReport(id, payload, len)
	default:
		fmt.Printf("fell through id: %d len: %d data: %v\n", id, len, payload)
		return nil
	}

}

type ReportInt interface {
	Handle()
}
