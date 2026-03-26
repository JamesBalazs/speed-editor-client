package hardware

const (
	JOGMODE_RELATIVE          = iota // Relative
	JOGMODE_ABSOLUTE                 // Absolute, where 0 is the position when the mode was set. -4096 -> 4096 = 180deg
	JOGMODE_RELATIVE_2               // Used for faster scrolling in later versions I think. Reports the same as JOGMODE_RELATIVE but we should eg double the speed in software
	JOGMODE_ABSOLUTE_DEADZONE        // Absolute, with a deadzone around 0

	JOG_MAX = 4096
	JOG_MIN = -4096
)
