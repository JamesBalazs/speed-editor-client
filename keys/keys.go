package keys

var (
	NullKey = Key{Name: NONE, Id: ID_NONE, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_NONE, SubText: SUBTEXT_NONE, Row: -1, Col: -1, Width: 1}
	// Keys is used when a deterministic ordering is required when iterating over keys, since
	// iterating the values of a map will not always return keys in the same order.
	keys = []Key{
		// Top-left group (rows 0-1, cols 0-2)
		{Name: SMART_INSRT, Id: ID_SMART_INSRT, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SMART_INSRT, SubText: SUBTEXT_SMART_INSRT, Row: 0, Col: 0, Width: 1},
		{Name: APPND, Id: ID_APPND, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_APPND, SubText: SUBTEXT_APPND, Row: 0, Col: 1, Width: 1},
		{Name: RIPL_OWR, Id: ID_RIPL_OWR, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_RIPL_OWR, SubText: SUBTEXT_RIPL_OWR, Row: 0, Col: 2, Width: 1},
		{Name: CLOSE_UP, Id: ID_CLOSE_UP, Led: LED_CLOSE_UP, JogLed: LED_NONE, Text: TEXT_CLOSE_UP, SubText: SUBTEXT_CLOSE_UP, Row: 1, Col: 0, Width: 1},
		{Name: PLACE_ON_TOP, Id: ID_PLACE_ON_TOP, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_PLACE_ON_TOP, SubText: SUBTEXT_PLACE_ON_TOP, Row: 1, Col: 1, Width: 1},
		{Name: SRC_OWR, Id: ID_SRC_OWR, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SRC_OWR, SubText: SUBTEXT_SRC_OWR, Row: 1, Col: 2, Width: 1},

		// IN/OUT row (wide keys, rows 2)
		{Name: IN, Id: ID_IN, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_IN, SubText: SUBTEXT_IN, Row: 2, Col: 0, Width: 1.5},
		{Name: OUT, Id: ID_OUT, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_OUT, SubText: SUBTEXT_OUT, Row: 2, Col: 1.5, Width: 1.5},

		// Bottom-left group (rows 3-5, cols 0-2)
		{Name: TRIM_IN, Id: ID_TRIM_IN, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_TRIM_IN, SubText: SUBTEXT_TRIM_IN, Row: 3, Col: 0, Width: 1},
		{Name: TRIM_OUT, Id: ID_TRIM_OUT, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_TRIM_OUT, SubText: SUBTEXT_TRIM_OUT, Row: 3, Col: 1, Width: 1},
		{Name: ROLL, Id: ID_ROLL, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_ROLL, SubText: SUBTEXT_ROLL, Row: 3, Col: 2, Width: 1},
		{Name: SLIP_SRC, Id: ID_SLIP_SRC, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SLIP_SRC, SubText: SUBTEXT_SLIP_SRC, Row: 4, Col: 0, Width: 1},
		{Name: SLIP_DEST, Id: ID_SLIP_DEST, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SLIP_DEST, SubText: SUBTEXT_SLIP_DEST, Row: 4, Col: 1, Width: 1},
		{Name: TRANS_DUR, Id: ID_TRANS_DUR, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_TRANS_DUR, SubText: SUBTEXT_TRANS_DUR, Row: 4, Col: 2, Width: 1},
		{Name: CUT, Id: ID_CUT, Led: LED_CUT, JogLed: LED_NONE, Text: TEXT_CUT, SubText: SUBTEXT_CUT, Row: 5, Col: 0, Width: 1},
		{Name: DIS, Id: ID_DIS, Led: LED_DIS, JogLed: LED_NONE, Text: TEXT_DIS, SubText: SUBTEXT_DIS, Row: 5, Col: 1, Width: 1},
		{Name: SMTH_CUT, Id: ID_SMTH_CUT, Led: LED_SMTH_CUT, JogLed: LED_NONE, Text: TEXT_SMTH_CUT, SubText: SUBTEXT_SMTH_CUT, Row: 5, Col: 2, Width: 1},

		// Top-centre group (rows 0-1, cols 3-6)
		{Name: ESC, Id: ID_ESC, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_ESC, SubText: SUBTEXT_ESC, Row: 0, Col: 3, Width: 1},
		{Name: SYNC_BIN, Id: ID_SYNC_BIN, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SYNC_BIN, SubText: SUBTEXT_SYNC_BIN, Row: 0, Col: 4, Width: 1},
		{Name: AUDIO_LEVEL, Id: ID_AUDIO_LEVEL, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_AUDIO_LEVEL, SubText: SUBTEXT_AUDIO_LEVEL, Row: 0, Col: 5, Width: 1},
		{Name: FULL_VIEW, Id: ID_FULL_VIEW, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_FULL_VIEW, SubText: SUBTEXT_FULL_VIEW, Row: 0, Col: 6, Width: 1},
		{Name: TRANS, Id: ID_TRANS, Led: LED_TRANS, JogLed: LED_NONE, Text: TEXT_TRANS, SubText: SUBTEXT_TRANS, Row: 1, Col: 3, Width: 1},
		{Name: SPLIT, Id: ID_SPLIT, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SPLIT, SubText: SUBTEXT_SPLIT, Row: 1, Col: 4, Width: 1},
		{Name: SNAP, Id: ID_SNAP, Led: LED_SNAP, JogLed: LED_NONE, Text: TEXT_SNAP, SubText: SUBTEXT_SNAP, Row: 1, Col: 5, Width: 1},
		{Name: RIPL_DEL, Id: ID_RIPL_DEL, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_RIPL_DEL, SubText: SUBTEXT_RIPL_DEL, Row: 1, Col: 6, Width: 1},

		// CAM grid (rows 2-4, cols 3-6)
		{Name: CAM7, Id: ID_CAM7, Led: LED_CAM7, JogLed: LED_NONE, Text: TEXT_CAM7, SubText: SUBTEXT_CAM7, Row: 2, Col: 3, Width: 1},
		{Name: CAM8, Id: ID_CAM8, Led: LED_CAM8, JogLed: LED_NONE, Text: TEXT_CAM8, SubText: SUBTEXT_CAM8, Row: 2, Col: 4, Width: 1},
		{Name: CAM9, Id: ID_CAM9, Led: LED_CAM9, JogLed: LED_NONE, Text: TEXT_CAM9, SubText: SUBTEXT_CAM9, Row: 2, Col: 5, Width: 1},
		{Name: LIVE_OWR, Id: ID_LIVE_OWR, Led: LED_LIVE_OWR, JogLed: LED_NONE, Text: TEXT_LIVE_OWR, SubText: SUBTEXT_LIVE_OWR, Row: 2, Col: 6, Width: 1},
		{Name: CAM4, Id: ID_CAM4, Led: LED_CAM4, JogLed: LED_NONE, Text: TEXT_CAM4, SubText: SUBTEXT_CAM4, Row: 3, Col: 3, Width: 1},
		{Name: CAM5, Id: ID_CAM5, Led: LED_CAM5, JogLed: LED_NONE, Text: TEXT_CAM5, SubText: SUBTEXT_CAM5, Row: 3, Col: 4, Width: 1},
		{Name: CAM6, Id: ID_CAM6, Led: LED_CAM6, JogLed: LED_NONE, Text: TEXT_CAM6, SubText: SUBTEXT_CAM6, Row: 3, Col: 5, Width: 1},
		{Name: VIDEO_ONLY, Id: ID_VIDEO_ONLY, Led: LED_VIDEO_ONLY, JogLed: LED_NONE, Text: TEXT_VIDEO_ONLY, SubText: SUBTEXT_VIDEO_ONLY, Row: 3, Col: 6, Width: 1},
		{Name: CAM1, Id: ID_CAM1, Led: LED_CAM1, JogLed: LED_NONE, Text: TEXT_CAM1, SubText: SUBTEXT_CAM1, Row: 4, Col: 3, Width: 1},
		{Name: CAM2, Id: ID_CAM2, Led: LED_CAM2, JogLed: LED_NONE, Text: TEXT_CAM2, SubText: SUBTEXT_CAM2, Row: 4, Col: 4, Width: 1},
		{Name: CAM3, Id: ID_CAM3, Led: LED_CAM3, JogLed: LED_NONE, Text: TEXT_CAM3, SubText: SUBTEXT_CAM3, Row: 4, Col: 5, Width: 1},
		{Name: AUDIO_ONLY, Id: ID_AUDIO_ONLY, Led: LED_AUDIO_ONLY, JogLed: LED_NONE, Text: TEXT_AUDIO_ONLY, SubText: SUBTEXT_AUDIO_ONLY, Row: 4, Col: 6, Width: 1},

		// STOP/PLAY spans the bottom of the CAM grid
		{Name: STOP_PLAY, Id: ID_STOP_PLAY, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_STOP_PLAY, SubText: SUBTEXT_STOP_PLAY, Row: 5, Col: 3, Width: 4},

		// Top-right group
		{Name: SOURCE, Id: ID_SOURCE, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_SOURCE, SubText: SUBTEXT_SOURCE, Row: 0, Col: 7, Width: 1.5},
		{Name: TIMELINE, Id: ID_TIMELINE, Led: LED_NONE, JogLed: LED_NONE, Text: TEXT_TIMELINE, SubText: SUBTEXT_TIMELINE, Row: 0, Col: 8.5, Width: 1.5},

		// Jog mode buttons
		{Name: SHTL, Id: ID_SHTL, Led: LED_NONE, JogLed: LED_SHTL, Text: TEXT_SHTL, SubText: SUBTEXT_SHTL, Row: 1, Col: 7, Width: 1},
		{Name: JOG, Id: ID_JOG, Led: LED_NONE, JogLed: LED_JOG, Text: TEXT_JOG, SubText: SUBTEXT_JOG, Row: 1, Col: 8, Width: 1},
		{Name: SCRL, Id: ID_SCRL, Led: LED_NONE, JogLed: LED_SCRL, Text: TEXT_SCRL, SubText: SUBTEXT_SCRL, Row: 1, Col: 9, Width: 1},
	}
)

type Key struct {
	Name string
	Id   uint16

	Led    uint32
	JogLed uint8

	Text    string
	SubText string

	Row   int
	Col   float32
	Width float32
}

// Get returns a new slice of all Keys each time it is called.
func Get() []Key {
	keysCopy := make([]Key, len(keys))
	copy(keys, keysCopy)

	return keysCopy
}

// ByName returns a map of keys, for constant time lookup by their Name.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByName() map[string]Key {
	keyIndex := make(map[string]Key, 43)

	for _, key := range keys {
		keyIndex[key.Name] = key
	}

	return keyIndex
}

// ById returns a map of keys, for constant time lookup by their Id.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ById() map[uint16]Key {
	keyIndex := make(map[uint16]Key, 43)

	for _, key := range keys {
		keyIndex[key.Id] = key
	}

	return keyIndex
}

// ByLedId returns a map of keys, for constant time lookup by their LedId.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByLedId() map[uint32]Key {
	keyIndex := make(map[uint32]Key, 43)

	for _, key := range keys {
		keyIndex[key.Led] = key
	}

	return keyIndex
}

// ByJogLedId returns a map of keys, for constant time lookup by their JogLedId.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByJogLedId() map[uint8]Key {
	keyIndex := make(map[uint8]Key, 43)

	for _, key := range keys {
		keyIndex[key.JogLed] = key
	}

	return keyIndex
}

// ByText returns a map of keys, for constant time lookup by their Text.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByText() map[string]Key {
	keyIndex := make(map[string]Key, 43)

	for _, key := range keys {
		keyIndex[key.Text] = key
	}

	return keyIndex
}

// BySubText returns a map of keys, for constant time lookup by their SubText.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func BySubText() map[string]Key {
	keyIndex := make(map[string]Key, 43)

	for _, key := range keys {
		keyIndex[key.SubText] = key
	}

	return keyIndex
}

// ByCol returns a map of keys, for constant time lookup by their X
// and Y positions respectively.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByCol() map[float32]map[int]Key {
	keyIndex := make(map[float32]map[int]Key, 10)

	for _, key := range keys {
		if keyIndex[key.Col] == nil {
			keyIndex[key.Col] = map[int]Key{}
		}
		keyIndex[key.Col][key.Row] = key
	}

	return keyIndex
}

// ByRow returns a map of keys, for constant time lookup by their Y
// and X positions respectively.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByRow() map[int]map[float32]Key {
	keyIndex := make(map[int]map[float32]Key, 6)

	for _, key := range keys {
		if keyIndex[key.Row] == nil {
			keyIndex[key.Row] = map[float32]Key{}
		}
		keyIndex[key.Row][key.Col] = key
	}

	return keyIndex
}
