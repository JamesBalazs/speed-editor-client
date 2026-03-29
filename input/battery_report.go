package input

import (
	"fmt"
)

func NewBatteryReport(id byte, payload []byte) (BatteryReport, error) {
	if id != ReportBattery {
		return BatteryReport{}, fmt.Errorf("malformed battery stats report id: %v payload: %v", id, payload)
	}

	return BatteryReport{
		Id:       id,
		Charging: payload[0] == 1,
		Battery:  float32(payload[1]) / 255,
	}, nil
}

type BatteryReport struct {
	Id       uint8
	Charging bool
	Battery  float32
}
