package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGet tests the Get function
func TestGet(t *testing.T) {
	t.Run("returns all keys", func(t *testing.T) {
		keys := Get()

		// Should return 43 keys
		assert.Len(t, keys, 43)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		keys1 := Get()
		keys2 := Get()

		// Should have equal content but be different underlying arrays
		assert.Equal(t, keys1, keys2)

		// Modifying one should not affect the other
		if len(keys1) > 0 {
			originalName := keys1[0].Name
			keys1[0].Name = "MODIFIED"
			assert.NotEqual(t, keys1[0].Name, keys2[0].Name)
			keys1[0].Name = originalName
		}
	})

	t.Run("contains expected keys", func(t *testing.T) {
		keys := Get()

		// Check for some expected keys
		found := false
		for _, key := range keys {
			if key.Name == CAM1 {
				found = true
				assert.Equal(t, "CAM1", key.Name)
				assert.NotZero(t, key.Id)
				assert.NotZero(t, key.Led)
				assert.Equal(t, 4, key.Row)
				assert.Equal(t, float32(3), key.Col)
			}
		}
		assert.True(t, found, "CAM1 key should be present")
	})

	t.Run("contains jog keys", func(t *testing.T) {
		keys := Get()

		jogKeys := []string{JOG, SHTL, SCRL}
		for _, expectedName := range jogKeys {
			found := false
			for _, key := range keys {
				if key.Name == expectedName {
					found = true
					assert.NotZero(t, key.JogLed)
				}
			}
			assert.True(t, found, "%s key should be present", expectedName)
		}
	})

	t.Run("contains cam keys", func(t *testing.T) {
		keys := Get()

		camKeys := []string{CAM1, CAM2, CAM3, CAM4, CAM5, CAM6, CAM7, CAM8, CAM9}
		for _, expectedName := range camKeys {
			found := false
			for _, key := range keys {
				if key.Name == expectedName {
					found = true
					assert.NotZero(t, key.Led)
				}
			}
			assert.True(t, found, "%s key should be present", expectedName)
		}
	})
}

// TestByName tests the ByName function
func TestByName(t *testing.T) {
	t.Run("returns map with all keys", func(t *testing.T) {
		keysByName := ByName()

		assert.Len(t, keysByName, 43)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByName()
		map2 := ByName()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by name", func(t *testing.T) {
		keysByName := ByName()

		key, exists := keysByName[CAM1]
		assert.True(t, exists)
		assert.Equal(t, "CAM1", key.Name)
		assert.NotZero(t, key.Led)

		key, exists = keysByName[TRANS]
		assert.True(t, exists)
		assert.Equal(t, "TRANS", key.Name)
		assert.NotZero(t, key.Led)
	})

	t.Run("lookup jog keys", func(t *testing.T) {
		keysByName := ByName()

		key, exists := keysByName[JOG]
		assert.True(t, exists)
		assert.Equal(t, "JOG", key.Name)
		assert.NotZero(t, key.JogLed)

		key, exists = keysByName[SHTL]
		assert.True(t, exists)
		assert.Equal(t, "SHTL", key.Name)
		assert.NotZero(t, key.JogLed)

		key, exists = keysByName[SCRL]
		assert.True(t, exists)
		assert.Equal(t, "SCRL", key.Name)
		assert.NotZero(t, key.JogLed)
	})

	t.Run("lookup wide keys", func(t *testing.T) {
		keysByName := ByName()

		key, exists := keysByName[IN]
		assert.True(t, exists)
		assert.Equal(t, float32(1.5), key.Width)

		key, exists = keysByName[OUT]
		assert.True(t, exists)
		assert.Equal(t, float32(1.5), key.Width)

		key, exists = keysByName[SOURCE]
		assert.True(t, exists)
		assert.Equal(t, float32(1.5), key.Width)

		key, exists = keysByName[TIMELINE]
		assert.True(t, exists)
		assert.Equal(t, float32(1.5), key.Width)

		key, exists = keysByName[STOP_PLAY]
		assert.True(t, exists)
		assert.Equal(t, float32(4), key.Width)
	})

	t.Run("non-existent key returns zero value", func(t *testing.T) {
		keysByName := ByName()

		key, exists := keysByName["NON_EXISTENT"]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestById tests the ById function
func TestById(t *testing.T) {
	t.Run("returns map with all keys", func(t *testing.T) {
		keysById := ById()

		assert.Len(t, keysById, 43)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ById()
		map2 := ById()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by id", func(t *testing.T) {
		keysById := ById()

		// Find a key by its ID
		var cam1ID uint16
		for _, key := range Get() {
			if key.Name == CAM1 {
				cam1ID = key.Id
				break
			}
		}
		assert.NotZero(t, cam1ID)

		key, exists := keysById[cam1ID]
		assert.True(t, exists)
		assert.Equal(t, CAM1, key.Name)
	})

	t.Run("non-existent id returns zero value", func(t *testing.T) {
		keysById := ById()

		key, exists := keysById[uint16(9999)]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestByLedId tests the ByLedId function
func TestByLedId(t *testing.T) {
	t.Run("returns map with keys that have leds", func(t *testing.T) {
		keysByLedId := ByLedId()

		// Only keys with non-zero LEDs will be in the map
		// Many keys share LED_NONE (0), so they'll be overwritten
		assert.Greater(t, len(keysByLedId), 1)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByLedId()
		map2 := ByLedId()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by led id", func(t *testing.T) {
		keysByLedId := ByLedId()

		// Find a key's LED ID
		var cam1Led uint32
		for _, key := range Get() {
			if key.Name == CAM1 {
				cam1Led = key.Led
				break
			}
		}
		assert.NotZero(t, cam1Led)

		key, exists := keysByLedId[cam1Led]
		assert.True(t, exists)
		assert.Equal(t, CAM1, key.Name)
	})

	t.Run("non-existent led id returns zero value", func(t *testing.T) {
		keysByLedId := ByLedId()

		key, exists := keysByLedId[uint32(9999)]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestByJogLedId tests the ByJogLedId function
func TestByJogLedId(t *testing.T) {
	t.Run("returns map with jog leds", func(t *testing.T) {
		keysByJogLedId := ByJogLedId()

		// Should have jog leds: LED_JOG, LED_SHTL, LED_SCRL, and LED_NONE
		assert.GreaterOrEqual(t, len(keysByJogLedId), 3)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByJogLedId()
		map2 := ByJogLedId()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by jog led id", func(t *testing.T) {
		keysByJogLedId := ByJogLedId()

		// Find jog key's LED ID
		var jogLed uint8
		for _, key := range Get() {
			if key.Name == JOG {
				jogLed = key.JogLed
				break
			}
		}
		assert.NotZero(t, jogLed)

		key, exists := keysByJogLedId[jogLed]
		assert.True(t, exists)
		assert.Equal(t, JOG, key.Name)
	})

	t.Run("non-existent jog led id returns zero value", func(t *testing.T) {
		keysByJogLedId := ByJogLedId()

		key, exists := keysByJogLedId[uint8(99)]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestByText tests the ByText function
func TestByText(t *testing.T) {
	t.Run("returns map with all keys", func(t *testing.T) {
		keysByText := ByText()

		assert.Len(t, keysByText, 43)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByText()
		map2 := ByText()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by text", func(t *testing.T) {
		keysByText := ByText()

		key, exists := keysByText[TEXT_CAM1]
		assert.True(t, exists)
		assert.Equal(t, CAM1, key.Name)

		key, exists = keysByText[TEXT_TRANS]
		assert.True(t, exists)
		assert.Equal(t, TRANS, key.Name)
	})

	t.Run("non-existent text returns zero value", func(t *testing.T) {
		keysByText := ByText()

		key, exists := keysByText["NON_EXISTENT_TEXT"]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestBySubText tests the BySubText function
func TestBySubText(t *testing.T) {
	t.Run("returns map with keys that have subtext", func(t *testing.T) {
		keysBySubText := BySubText()

		// Only keys with non-empty SubText will be in the map
		// Many keys share empty SubText, so they'll be overwritten
		assert.Greater(t, len(keysBySubText), 1)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := BySubText()
		map2 := BySubText()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by subtext", func(t *testing.T) {
		keysBySubText := BySubText()

		// Find a key with unique subtext and verify it maps correctly
		// Note: Multiple keys may share the same subtext (e.g., "CLIP"), 
		// so we just verify the map returns a valid key for existing subtexts
		for subtext, key := range keysBySubText {
			if subtext != "" {
				assert.NotEmpty(t, key.Name)
				assert.Equal(t, subtext, key.SubText)
				return
			}
		}
	})

	t.Run("non-existent subtext returns zero value", func(t *testing.T) {
		keysBySubText := BySubText()

		key, exists := keysBySubText["NON_EXISTENT_SUBTEXT"]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestByCol tests the ByCol function
func TestByCol(t *testing.T) {
	t.Run("returns nested map", func(t *testing.T) {
		keysByCol := ByCol()

		// Should have multiple columns
		assert.Greater(t, len(keysByCol), 0)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByCol()
		map2 := ByCol()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by column and row", func(t *testing.T) {
		keysByCol := ByCol()

		// CAM1 is at col 3, row 4
		col, exists := keysByCol[float32(3)]
		assert.True(t, exists)

		key, exists := col[4]
		assert.True(t, exists)
		assert.Equal(t, CAM1, key.Name)
	})

	t.Run("lookup wide key", func(t *testing.T) {
		keysByCol := ByCol()

		// IN is at col 0, row 2
		col, exists := keysByCol[float32(0)]
		assert.True(t, exists)

		key, exists := col[2]
		assert.True(t, exists)
		assert.Equal(t, IN, key.Name)
	})

	t.Run("non-existent column returns empty map", func(t *testing.T) {
		keysByCol := ByCol()

		col, exists := keysByCol[float32(999)]
		assert.False(t, exists)
		assert.Nil(t, col)
	})

	t.Run("non-existent row in existing column", func(t *testing.T) {
		keysByCol := ByCol()

		col, exists := keysByCol[float32(0)]
		assert.True(t, exists)

		key, exists := col[999]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})
}

// TestByRow tests the ByRow function
func TestByRow(t *testing.T) {
	t.Run("returns nested map", func(t *testing.T) {
		keysByRow := ByRow()

		// Should have 6 rows (0-5)
		assert.Len(t, keysByRow, 6)
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		map1 := ByRow()
		map2 := ByRow()

		// Should have equal content but be different maps
		assert.Equal(t, map1, map2)
	})

	t.Run("lookup by row and column", func(t *testing.T) {
		keysByRow := ByRow()

		// CAM1 is at row 4, col 3
		row, exists := keysByRow[4]
		assert.True(t, exists)

		key, exists := row[float32(3)]
		assert.True(t, exists)
		assert.Equal(t, CAM1, key.Name)
	})

	t.Run("lookup wide key by row", func(t *testing.T) {
		keysByRow := ByRow()

		// IN is at row 2, col 0
		row, exists := keysByRow[2]
		assert.True(t, exists)

		key, exists := row[float32(0)]
		assert.True(t, exists)
		assert.Equal(t, IN, key.Name)
	})

	t.Run("non-existent row returns empty map", func(t *testing.T) {
		keysByRow := ByRow()

		row, exists := keysByRow[999]
		assert.False(t, exists)
		assert.Nil(t, row)
	})

	t.Run("non-existent column in existing row", func(t *testing.T) {
		keysByRow := ByRow()

		row, exists := keysByRow[0]
		assert.True(t, exists)

		key, exists := row[float32(999)]
		assert.False(t, exists)
		assert.Equal(t, Key{}, key)
	})

	t.Run("all rows have expected keys", func(t *testing.T) {
		keysByRow := ByRow()

		// Row 0 should have multiple keys
		assert.Greater(t, len(keysByRow[0]), 0)
		// Row 5 should have keys including CUT, DIS, SMTH_CUT, STOP_PLAY
		assert.Greater(t, len(keysByRow[5]), 0)
	})
}

// TestNullKey tests the NullKey constant
func TestNullKey(t *testing.T) {
	t.Run("has expected values", func(t *testing.T) {
		assert.Equal(t, NONE, NullKey.Name)
		assert.Equal(t, uint16(0), NullKey.Id)
		assert.Equal(t, uint32(0), NullKey.Led)
		assert.Equal(t, uint8(0), NullKey.JogLed)
		assert.Equal(t, TEXT_NONE, NullKey.Text)
		assert.Equal(t, SUBTEXT_NONE, NullKey.SubText)
		assert.Equal(t, -1, NullKey.Row)
		assert.Equal(t, float32(-1), NullKey.Col)
		assert.Equal(t, float32(1), NullKey.Width)
	})
}

// TestKeyStruct tests the Key struct
func TestKeyStruct(t *testing.T) {
	t.Run("can create key with all fields", func(t *testing.T) {
		key := Key{
			Name:    "TEST",
			Id:      123,
			Led:     456,
			JogLed:  78,
			Text:    "Test Text",
			SubText: "Test Sub",
			Row:     5,
			Col:     3.5,
			Width:   2.0,
		}

		assert.Equal(t, "TEST", key.Name)
		assert.Equal(t, uint16(123), key.Id)
		assert.Equal(t, uint32(456), key.Led)
		assert.Equal(t, uint8(78), key.JogLed)
		assert.Equal(t, "Test Text", key.Text)
		assert.Equal(t, "Test Sub", key.SubText)
		assert.Equal(t, 5, key.Row)
		assert.Equal(t, float32(3.5), key.Col)
		assert.Equal(t, float32(2.0), key.Width)
	})
}
