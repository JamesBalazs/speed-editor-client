package inputReport

import (
	"encoding/binary"
	"log"
)

func NewJogReport(id byte, payload []byte) JogReport {
	if id != ReportJog {
		log.Fatalf("malformed jog input report id: %v payload: %", id, payload)
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
