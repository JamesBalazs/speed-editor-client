package inputReport

import (
	"fmt"
	"log"
)

func NewBatteryReport(id byte, payload []byte, length int) ReportInt {
	if id != ReportBattery || length != 2 || len(payload) != length {
		log.Fatalf("malformed battery stats report id: %v payload: %v len: %d", id, payload, length)
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

func (report BatteryReport) Handle() {
	fmt.Printf("got battery %v %f\n", report.Charging, report.Battery)

	return // TODO implement handler
}
