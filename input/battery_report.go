package input

import (
	"log"
)

func NewBatteryReport(id byte, payload []byte) BatteryReport {
	if id != ReportBattery {
		log.Fatalf("malformed battery stats report id: %v payload: %v", id, payload)
	}

	return BatteryReport{
		Id:       id,
		Charging: payload[0] == 1,
		Battery:  float32(payload[1]) / 255,
	}
}

type BatteryReport struct {
	Id       uint8
	Charging bool
	Battery  float32
}
