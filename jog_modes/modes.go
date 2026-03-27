package jogModes

const (
	ID_RELATIVE          = iota // Relative
	ID_ABSOLUTE                 // Absolute, where 0 is the position when the mode was set. -4096 -> 4096 = 180deg
	ID_RELATIVE_2               // Used for faster scrolling in later versions I think. Reports the same as ID_RELATIVE but we should eg double the speed in software
	ID_ABSOLUTE_DEADZONE        // Absolute, with a deadzone around 0

	RELATIVE          = "RELATIVE"
	ABSOLUTE          = "ABSOLUTE"
	RELATIVE_2        = "RELATIVE_2"
	ABSOLUTE_DEADZONE = "ABSOLUTE_DEADZONE"

	ABSOLUTE_MAX = 4096
	ABSOLUTE_MIN = -4096
)

type Mode struct {
	Id   int
	Name string
}

var modes = []Mode{
	{Id: ID_RELATIVE, Name: RELATIVE},
	{Id: ID_ABSOLUTE, Name: ABSOLUTE},
	{Id: ID_RELATIVE_2, Name: RELATIVE_2},
	{Id: ID_ABSOLUTE_DEADZONE, Name: ABSOLUTE_DEADZONE},
}

// Get returns a new slice of all jog Modes each time it is called.
func Get() []Mode {
	modesCopy := make([]Mode, len(modes))
	copy(modes, modesCopy)

	return modesCopy
}

// ByName returns a map of modes, for constant time lookup by their Name.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ByName() map[string]Mode {
	modeIndex := make(map[string]Mode, 43)

	for _, mode := range modes {
		modeIndex[mode.Name] = mode
	}

	return modeIndex
}

// ById returns a map of modes, for constant time lookup by their Id.
// A new copy of the map is returned each time, so when the consumer modifies
// the map it doesn't modify the underlying data.
func ById() map[int]Mode {
	modeIndex := make(map[int]Mode, 43)

	for _, mode := range modes {
		modeIndex[mode.Id] = mode
	}

	return modeIndex
}
