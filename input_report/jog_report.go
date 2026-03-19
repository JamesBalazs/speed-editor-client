package inputReport

import (
	"encoding/binary"
	"fmt"
	"log"
)

func NewJogReport(id byte, payload []byte, length int) ReportInt {
	if id != ReportJog || length != 6 || len(payload) != length {
		log.Fatalf("malformed jog input report id: %v payload: %v len: %d", id, payload, length)
	}

	return JogReport{
		Id:      id,
		Mode:    payload[0],
		Value:   int32(binary.LittleEndian.Uint32(payload[1:5])),
		Unknown: payload[5],
	}
}

type JogReport struct {
	Id      uint8
	Mode    uint8
	Value   int32
	Unknown uint8
}

func (report JogReport) Handle() {
	fmt.Printf("got jog mode %d pos %d\n", report.Mode, report.Value)

	return // TODO implement handler
}
