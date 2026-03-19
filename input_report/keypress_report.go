package inputReport

import (
	"encoding/binary"
	"fmt"
	"log"
)

func NewKeyPressReport(id byte, payload []byte, length int) ReportInt {
	if id != ReportKeyPress || length != 12 || len(payload) != length {
		log.Fatalf("malformed keypress input report id: %v payload: %v len: %d", id, payload, length)
	}

	keys := make([]int16, 6)
	for i := 0; i < length; i += 2 {
		keys[i/2] = int16(binary.LittleEndian.Uint16(payload[i : i+2]))
	}

	return KeyPressReport{
		Id:   id,
		Keys: keys,
	}
}

type KeyPressReport struct {
	Id   uint8
	Keys []int16
}

func (report KeyPressReport) Handle() {
	fmt.Printf("got keypresses %v\n", report.Keys)

	return // TODO implement handler
}
