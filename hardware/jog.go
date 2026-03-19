package hardware

const (
	JOGMODE_RELATIVE          = "RELATIVE"   // Relative
	JOGMODE_ABSOLUTE          = "ABSOLUTE"   // Absolute, where 0 is the position when the mode was set. -4096 -> 4096 = 180deg
	JOGMODE_RELATIVE_2        = "RELATIVE2"  // Used for faster scrolling in later versions I think. Reports the same as JOGMODE_RELATIVE but we should eg double the speed in software
	JOGMODE_ABSOLUTE_DEADZONE = "ABSOLUTE_0" // Absolute, with a deadzone around 0
)

var JogMode = map[string]byte{
	JOGMODE_RELATIVE:          0,
	JOGMODE_ABSOLUTE:          1,
	JOGMODE_RELATIVE_2:        2,
	JOGMODE_ABSOLUTE_DEADZONE: 3,
}
