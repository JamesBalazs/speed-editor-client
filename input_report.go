package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

const (
	ReportJog          = 3
	ReportKeyPress     = 4
	ReportBatteryStats = 7
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
	default:
		fmt.Printf("fell through id: %d len: %d data: %v\n", id, len, payload)
		return nil
	}

}

type ReportInt interface {
	Handle()
}

func NewJogReport(id byte, payload []byte, length int) ReportInt {
	if id != 3 || length != 6 || len(payload) != length {
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

func NewKeyPressReport(id byte, payload []byte, length int) ReportInt {
	if id != 4 || length != 12 || len(payload) != length {
		log.Fatalf("malformed keypress input report id: %v payload: %v len: %d", id, payload, length)
	}

	keys := make([]int16, 6)
	for i := 0; i < length; i += 2 {
		keys[i/2] = int16(binary.LittleEndian.Uint16(payload[i:i+2]))
	}

	return KeyPressReport{
		Id:      id,
		Keys:    keys,
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
