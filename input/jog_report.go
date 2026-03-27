package input

import (
	"encoding/binary"
	"fmt"

	jogModes "github.com/JamesBalazs/speed-editor-client/jog_modes"
)

var modesById = jogModes.ById()

func NewJogReport(id byte, payload []byte) JogReport {
	if id != ReportJog {
		fmt.Printf("malformed jog input report id: %v payload: %v\n", id, payload)
	}

	return JogReport{
		Id:      id,
		Mode:    modesById[int(payload[0])],
		Value:   int32(binary.LittleEndian.Uint32(payload[1:5])),
		Unknown: payload[5],
	}
}

type JogReport struct {
	Id      uint8
	Mode    jogModes.Mode
	Value   int32
	Unknown uint8
}
