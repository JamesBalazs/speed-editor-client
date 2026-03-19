package keys

const (
	NONE = "NONE" // /

	SMART_INSRT  = "SMART_INSRT"  // SMART INSRT [CLIP]
	APPND        = "APPND"        // APPND [CLIP]
	RIPL_OWR     = "RIPL_OWR"     // RIPL O/WR
	CLOSE_UP     = "CLOSE_UP"     // CLOSE UP [YPOS]
	PLACE_ON_TOP = "PLACE_ON_TOP" // PLACE ON TOP
	SRC_OWR      = "SRC_OWR"      // SRC O/WR

	IN        = "IN"        // IN [CLR]
	OUT       = "OUT"       // OUT [CLR]
	TRIM_IN   = "TRIM_IN"   // TRIM IN
	TRIM_OUT  = "TRIM_OUT"  // TRIM OUT
	ROLL      = "ROLL"      // ROLL [SLIDE]
	SLIP_SRC  = "SLIP_SRC"  // SLIP SRC
	SLIP_DEST = "SLIP_DEST" // SLIP DEST
	TRANS_DUR = "TRANS_DUR" // TRANS DUR [SET]
	CUT       = "CUT"       // CUT
	DIS       = "DIS"       // DIS
	SMTH_CUT  = "SMTH_CUT"  // SMTH CUT

	SOURCE   = "SOURCE"   // SOURCE
	TIMELINE = "TIMELINE" // TIMELINE

	SHTL = "SHTL" // SHTL
	JOG  = "JOG"  // JOG
	SCRL = "SCRL" // SCRL

	ESC         = "ESC"         // ESC [UNDO]
	SYNC_BIN    = "SYNC_BIN"    // SYNC BIN
	AUDIO_LEVEL = "AUDIO_LEVEL" // AUDIO LEVEL [MARK]
	FULL_VIEW   = "FULL_VIEW"   // FULL VIEW [RVW]
	TRANS       = "TRANS"       // TRANS [TITLE]
	SPLIT       = "SPLIT"       // SPLIT [MOVE]
	SNAP        = "SNAP"        // SNAP = "SNAP" [=]
	RIPL_DEL    = "RIPL_DEL"    // RIPL DEL

	CAM1       = "CAM1"       // CAM1
	CAM2       = "CAM2"       // CAM2
	CAM3       = "CAM3"       // CAM3
	CAM4       = "CAM4"       // CAM4
	CAM5       = "CAM5"       // CAM5
	CAM6       = "CAM6"       // CAM6
	CAM7       = "CAM7"       // CAM7
	CAM8       = "CAM8"       // CAM8
	CAM9       = "CAM9"       // CAM9
	LIVE_OWR   = "LIVE_OWR"   // LIVE O/WR [RND]
	VIDEO_ONLY = "VIDEO_ONLY" // VIDEO ONLY
	AUDIO_ONLY = "AUDIO_ONLY" // AUDIO ONLY
	STOP_PLAY  = "STOP_PLAY"  // STOP/PLAY
)

var Keys = map[string]uint16{
	NONE: 0x00, // /

	SMART_INSRT:  0x01, // SMART INSRT [CLIP]
	APPND:        0x02, // APPND [CLIP]
	RIPL_OWR:     0x03, // RIPL O/WR
	CLOSE_UP:     0x04, // CLOSE UP [YPOS]
	PLACE_ON_TOP: 0x05, // PLACE ON TOP
	SRC_OWR:      0x06, // SRC O/WR

	IN:        0x07, // IN [CLR]
	OUT:       0x08, // OUT [CLR]
	TRIM_IN:   0x09, // TRIM IN
	TRIM_OUT:  0x0a, // TRIM OUT
	ROLL:      0x0b, // ROLL [SLIDE]
	SLIP_SRC:  0x0c, // SLIP SRC
	SLIP_DEST: 0x0d, // SLIP DEST
	TRANS_DUR: 0x0e, // TRANS DUR [SET]
	CUT:       0x0f, // CUT
	DIS:       0x10, // DIS
	SMTH_CUT:  0x11, // SMTH CUT

	SOURCE:   0x1a, // SOURCE
	TIMELINE: 0x1b, // TIMELINE

	SHTL: 0x1c, // SHTL
	JOG:  0x1d, // JOG
	SCRL: 0x1e, // SCRL

	ESC:         0x31, // ESC [UNDO]
	SYNC_BIN:    0x1f, // SYNC BIN
	AUDIO_LEVEL: 0x2c, // AUDIO LEVEL [MARK]
	FULL_VIEW:   0x2d, // FULL VIEW [RVW]
	TRANS:       0x22, // TRANS [TITLE]
	SPLIT:       0x2f, // SPLIT [MOVE]
	SNAP:        0x2e, // SNAP [=]
	RIPL_DEL:    0x2b, // RIPL DEL

	CAM1:       0x33, // CAM1
	CAM2:       0x34, // CAM2
	CAM3:       0x35, // CAM3
	CAM4:       0x36, // CAM4
	CAM5:       0x37, // CAM5
	CAM6:       0x38, // CAM6
	CAM7:       0x39, // CAM7
	CAM8:       0x3a, // CAM8
	CAM9:       0x3b, // CAM9
	LIVE_OWR:   0x30, // LIVE O/WR [RND]
	VIDEO_ONLY: 0x25, // VIDEO ONLY
	AUDIO_ONLY: 0x26, // AUDIO ONLY
	STOP_PLAY:  0x3c, // STOP/PLAY
}
