package keys

var (
	Leds = map[string]uint32{
		CLOSE_UP:   (1 << 0),
		CUT:        (1 << 1),
		DIS:        (1 << 2),
		SMTH_CUT:   (1 << 3),
		TRANS:      (1 << 4),
		SNAP:       (1 << 5),
		CAM7:       (1 << 6),
		CAM8:       (1 << 7),
		CAM9:       (1 << 8),
		LIVE_OWR:   (1 << 9),
		CAM4:       (1 << 10),
		CAM5:       (1 << 11),
		CAM6:       (1 << 12),
		VIDEO_ONLY: (1 << 13),
		CAM1:       (1 << 14),
		CAM2:       (1 << 15),
		CAM3:       (1 << 16),
		AUDIO_ONLY: (1 << 17),
	}

	JogLeds = map[string]byte{
		JOG:  (1 << 0),
		SHTL: (1 << 1),
		SCRL: (1 << 2),
	}
)
