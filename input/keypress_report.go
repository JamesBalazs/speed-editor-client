package input

import (
	"encoding/binary"
	"fmt"

	"github.com/JamesBalazs/speed-editor-client/keys"
)

func NewKeyPressReport(id byte, payload []byte) (KeyPressReport, error) {
	if id != ReportKeyPress {
		return KeyPressReport{}, fmt.Errorf("malformed keypress input report id: %v payload: %v", id, payload)
	}

	keys := []keys.Key{}
	for i := 0; i < len(payload); i += 2 {
		id := binary.LittleEndian.Uint16(payload[i : i+2])

		if key, found := keysById[id]; found {
			keys = append(keys, key)
		}
	}

	return KeyPressReport{
		Id:   id,
		Keys: keys,
	}, nil
}

type KeyPressReport struct {
	Id   uint8
	Keys []keys.Key
}
