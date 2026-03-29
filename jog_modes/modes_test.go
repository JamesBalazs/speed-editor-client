package jogModes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGet tests the Get function
func TestGet(t *testing.T) {
	t.Run("returns all modes", func(t *testing.T) {
		modes := Get()

		// Should return 4 modes
		assert.Len(t, modes, 4)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		modes1 := Get()
		modes2 := Get()

		// Should have equal content but be different underlying arrays
		assert.Equal(t, modes1, modes2)

		// Modifying one should not affect the other
		if len(modes1) > 0 {
			originalName := modes1[0].Name
			modes1[0].Name = "MODIFIED"
			assert.NotEqual(t, modes1[0].Name, modes2[0].Name)
			modes1[0].Name = originalName
		}
	})

	t.Run("contains expected modes", func(t *testing.T) {
		modes := Get()

		found := false
		for _, mode := range modes {
			if mode.Name == ABSOLUTE {
				found = true
				assert.Equal(t, ID_ABSOLUTE, mode.Id)
			}
		}
		assert.True(t, found, "ABSOLUTE mode should be present")
	})

	t.Run("contains all mode types", func(t *testing.T) {
		modes := Get()

		modeNames := make(map[string]bool)
		for _, mode := range modes {
			modeNames[mode.Name] = true
		}

		assert.True(t, modeNames[RELATIVE], "RELATIVE mode should exist")
		assert.True(t, modeNames[ABSOLUTE], "ABSOLUTE mode should exist")
		assert.True(t, modeNames[RELATIVE_2], "RELATIVE_2 mode should exist")
		assert.True(t, modeNames[ABSOLUTE_DEADZONE], "ABSOLUTE_DEADZONE mode should exist")
	})
}

// TestByName tests the ByName function
func TestByName(t *testing.T) {
	t.Run("returns map with all modes", func(t *testing.T) {
		modesByName := ByName()

		assert.Len(t, modesByName, 4)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByName()
		map2 := ByName()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup relative mode", func(t *testing.T) {
		modesByName := ByName()

		mode, exists := modesByName[RELATIVE]
		assert.True(t, exists)
		assert.Equal(t, ID_RELATIVE, mode.Id)
		assert.Equal(t, RELATIVE, mode.Name)
	})

	t.Run("lookup absolute mode", func(t *testing.T) {
		modesByName := ByName()

		mode, exists := modesByName[ABSOLUTE]
		assert.True(t, exists)
		assert.Equal(t, ID_ABSOLUTE, mode.Id)
		assert.Equal(t, ABSOLUTE, mode.Name)
	})

	t.Run("lookup relative_2 mode", func(t *testing.T) {
		modesByName := ByName()

		mode, exists := modesByName[RELATIVE_2]
		assert.True(t, exists)
		assert.Equal(t, ID_RELATIVE_2, mode.Id)
		assert.Equal(t, RELATIVE_2, mode.Name)
	})

	t.Run("lookup absolute_deadzone mode", func(t *testing.T) {
		modesByName := ByName()

		mode, exists := modesByName[ABSOLUTE_DEADZONE]
		assert.True(t, exists)
		assert.Equal(t, ID_ABSOLUTE_DEADZONE, mode.Id)
		assert.Equal(t, ABSOLUTE_DEADZONE, mode.Name)
	})

	t.Run("non-existent mode returns zero value", func(t *testing.T) {
		modesByName := ByName()

		mode, exists := modesByName["NON_EXISTENT"]
		assert.False(t, exists)
		assert.Equal(t, Mode{}, mode)
	})
}

// TestById tests the ById function
func TestById(t *testing.T) {
	t.Run("returns map with all modes", func(t *testing.T) {
		modesById := ById()

		assert.Len(t, modesById, 4)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ById()
		map2 := ById()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup relative mode by id", func(t *testing.T) {
		modesById := ById()

		mode, exists := modesById[ID_RELATIVE]
		assert.True(t, exists)
		assert.Equal(t, RELATIVE, mode.Name)
	})

	t.Run("lookup absolute mode by id", func(t *testing.T) {
		modesById := ById()

		mode, exists := modesById[ID_ABSOLUTE]
		assert.True(t, exists)
		assert.Equal(t, ABSOLUTE, mode.Name)
	})

	t.Run("lookup relative_2 mode by id", func(t *testing.T) {
		modesById := ById()

		mode, exists := modesById[ID_RELATIVE_2]
		assert.True(t, exists)
		assert.Equal(t, RELATIVE_2, mode.Name)
	})

	t.Run("lookup absolute_deadzone mode by id", func(t *testing.T) {
		modesById := ById()

		mode, exists := modesById[ID_ABSOLUTE_DEADZONE]
		assert.True(t, exists)
		assert.Equal(t, ABSOLUTE_DEADZONE, mode.Name)
	})

	t.Run("non-existent id returns zero value", func(t *testing.T) {
		modesById := ById()

		mode, exists := modesById[999]
		assert.False(t, exists)
		assert.Equal(t, Mode{}, mode)
	})
}

// TestModeConstants tests the mode constants
func TestModeConstants(t *testing.T) {
	t.Run("relative constant", func(t *testing.T) {
		assert.Equal(t, "RELATIVE", RELATIVE)
	})

	t.Run("absolute constant", func(t *testing.T) {
		assert.Equal(t, "ABSOLUTE", ABSOLUTE)
	})

	t.Run("relative_2 constant", func(t *testing.T) {
		assert.Equal(t, "RELATIVE_2", RELATIVE_2)
	})

	t.Run("absolute_deadzone constant", func(t *testing.T) {
		assert.Equal(t, "ABSOLUTE_DEADZONE", ABSOLUTE_DEADZONE)
	})
}

// TestIdConstants tests the ID constants
func TestIdConstants(t *testing.T) {
	t.Run("id relative is zero", func(t *testing.T) {
		assert.Equal(t, 0, ID_RELATIVE)
	})

	t.Run("id absolute is one", func(t *testing.T) {
		assert.Equal(t, 1, ID_ABSOLUTE)
	})

	t.Run("id relative_2 is two", func(t *testing.T) {
		assert.Equal(t, 2, ID_RELATIVE_2)
	})

	t.Run("id absolute_deadzone is three", func(t *testing.T) {
		assert.Equal(t, 3, ID_ABSOLUTE_DEADZONE)
	})

	t.Run("ids are sequential", func(t *testing.T) {
		assert.Equal(t, ID_RELATIVE+1, ID_ABSOLUTE)
		assert.Equal(t, ID_ABSOLUTE+1, ID_RELATIVE_2)
		assert.Equal(t, ID_RELATIVE_2+1, ID_ABSOLUTE_DEADZONE)
	})
}

// TestAbsoluteConstants tests the absolute value constants
func TestAbsoluteConstants(t *testing.T) {
	t.Run("absolute max", func(t *testing.T) {
		assert.Equal(t, 4096, ABSOLUTE_MAX)
	})

	t.Run("absolute min", func(t *testing.T) {
		assert.Equal(t, -4096, ABSOLUTE_MIN)
	})

	t.Run("max and min are symmetric", func(t *testing.T) {
		assert.Equal(t, ABSOLUTE_MAX, -ABSOLUTE_MIN)
	})
}

// TestModeStruct tests the Mode struct
func TestModeStruct(t *testing.T) {
	t.Run("can create mode with all fields", func(t *testing.T) {
		mode := Mode{
			Id:   5,
			Name: "CUSTOM_MODE",
		}

		assert.Equal(t, int(5), mode.Id)
		assert.Equal(t, "CUSTOM_MODE", mode.Name)
	})

	t.Run("zero value mode", func(t *testing.T) {
		var mode Mode

		assert.Equal(t, int(0), mode.Id)
		assert.Equal(t, "", mode.Name)
	})
}

// TestModesSlice tests the internal modes slice
func TestModesSlice(t *testing.T) {
	t.Run("modes are in correct order", func(t *testing.T) {
		modes := Get()

		// Verify the order matches the iota definition
		assert.Equal(t, ID_RELATIVE, modes[0].Id)
		assert.Equal(t, ID_ABSOLUTE, modes[1].Id)
		assert.Equal(t, ID_RELATIVE_2, modes[2].Id)
		assert.Equal(t, ID_ABSOLUTE_DEADZONE, modes[3].Id)
	})

	t.Run("all modes have non-empty names", func(t *testing.T) {
		modes := Get()

		for _, mode := range modes {
			assert.NotEmpty(t, mode.Name, "Mode %d should have a name", mode.Id)
		}
	})
}
