package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRol8 tests the rol8 function
func TestRol8(t *testing.T) {
	t.Run("rotate left by 8 bits", func(t *testing.T) {
		// 0x01 rotated left by 8 bits should move the high byte to low byte
		result := rol8(0x01)
		assert.Equal(t, uint64(0x0100000000000000), result)
	})

	t.Run("rotate 0x1234567890ABCDEF left by 8 bits", func(t *testing.T) {
		// Left rotate by 8: moves highest byte to lowest position
		result := rol8(0x1234567890ABCDEF)
		assert.Equal(t, uint64(0xEF1234567890ABCD), result)
	})

	t.Run("rotate 0xFF left by 8 bits", func(t *testing.T) {
		result := rol8(0xFF)
		assert.Equal(t, uint64(0xFF00000000000000), result)
	})

	t.Run("rotate 0x00 returns 0x00", func(t *testing.T) {
		result := rol8(0x00)
		assert.Equal(t, uint64(0x00), result)
	})

	t.Run("rotate max uint64", func(t *testing.T) {
		result := rol8(0xFFFFFFFFFFFFFFFF)
		assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), result)
	})
}

// TestRol8n tests the rol8n function
func TestRol8n(t *testing.T) {
	t.Run("rotate by 0 returns original value", func(t *testing.T) {
		value := uint64(0x1234567890ABCDEF)
		result := rol8n(value, 0)
		assert.Equal(t, value, result)
	})

	t.Run("rotate by 1 applies single rotation", func(t *testing.T) {
		value := uint64(0x1234567890ABCDEF)
		result := rol8n(value, 1)
		expected := rol8(value)
		assert.Equal(t, expected, result)
	})

	t.Run("rotate by 3 applies three rotations", func(t *testing.T) {
		value := uint64(0x1234567890ABCDEF)
		result := rol8n(value, 3)
		expected := rol8(rol8(rol8(value)))
		assert.Equal(t, expected, result)
	})

	t.Run("rotate by 8 (full cycle) returns original", func(t *testing.T) {
		// 8 rotations of 8 bits = 64 bits = full cycle
		value := uint64(0x1234567890ABCDEF)
		result := rol8n(value, 8)
		assert.Equal(t, value, result)
	})

	t.Run("rotate 0x01 by 4", func(t *testing.T) {
		result := rol8n(0x01, 4)
		assert.NotZero(t, result)
		assert.NotEqual(t, uint64(0x01), result)
	})
}

// TestCalculateChallengeResponse tests the challenge response calculation
func TestCalculateChallengeResponse(t *testing.T) {
	t.Run("known challenge-response pairs", func(t *testing.T) {
		testCases := []struct {
			challenge uint64
			response  uint64
		}{
			{0x0000000000000000, 0x3ae1206f97c10bc8},
			{0x0000000000000001, 0x2b9ab32bebf244c6},
			{0x0000000000000002, 0x20a4fab8df9adf0a},
			{0x0000000000000003, 0x6df72d1b40aef698},
			{0x0000000000000004, 0x72226f051e66ab94},
			{0x0000000000000005, 0x3831a3c6032d6a42},
			{0x0000000000000006, 0xfd7ff81881352889},
			{0x0000000000000007, 0x751bf623f42e0ade},
			{0x0000000000000008, 0x3ae1206f97c10bc0},
			{0x0000000000000009, 0x2392b32bebf244c6},
			{0x000000000000000A, 0x20acfab8df9adf0a},
			{0x000000000000000B, 0x6df7251340aef698},
			{0x000000000000000C, 0x72226f0d1666ab94},
			{0x000000000000000D, 0x3831a3c60b256242},
			{0x000000000000000E, 0xfd7ff818813d2889},
			{0x000000000000000F, 0x751bf623f42e02de},
			{0xFFFFFFFFFFFFFFFF, 0x61a3f6474ff236c6},
		}

		for _, tc := range testCases {
			t.Run(string(rune(tc.challenge)), func(t *testing.T) {
				result := CalculateChallengeResponse(tc.challenge)
				assert.Equal(t, tc.response, result, "Challenge 0x%016x should produce response 0x%016x", tc.challenge, tc.response)
			})
		}
	})

	t.Run("response is deterministic", func(t *testing.T) {
		challenges := []uint64{0, 1, 42, 100, 1000, 10000, 100000}
		for _, challenge := range challenges {
			t.Run(string(rune(challenge)), func(t *testing.T) {
				result1 := CalculateChallengeResponse(challenge)
				result2 := CalculateChallengeResponse(challenge)
				assert.Equal(t, result1, result2)
			})
		}
	})
}
