package input

import (
	"encoding/binary"
	"testing"

	"github.com/JamesBalazs/speed-editor-client/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewBatteryReport tests the NewBatteryReport function
func TestNewBatteryReport(t *testing.T) {
	t.Run("success charging", func(t *testing.T) {
		id := ReportBattery
		payload := []byte{0x01, 0x80} // Charging, ~50% battery

		report, err := NewBatteryReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		assert.True(t, report.Charging)
		assert.InDelta(t, 0.5, report.Battery, 0.01)
	})

	t.Run("success not charging", func(t *testing.T) {
		id := ReportBattery
		payload := []byte{0x00, 0xFF} // Not charging, 100% battery

		report, err := NewBatteryReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		assert.False(t, report.Charging)
		assert.InDelta(t, 1.0, report.Battery, 0.01)
	})

	t.Run("success zero battery", func(t *testing.T) {
		id := ReportBattery
		payload := []byte{0x01, 0x00} // Charging, 0% battery

		report, err := NewBatteryReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		assert.True(t, report.Charging)
		assert.Equal(t, float32(0), report.Battery)
	})

	t.Run("success partial battery", func(t *testing.T) {
		id := ReportBattery
		payload := []byte{0x00, 0x40} // Not charging, ~25% battery

		report, err := NewBatteryReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		assert.False(t, report.Charging)
		assert.InDelta(t, 0.25, report.Battery, 0.01)
	})

	t.Run("error wrong report id", func(t *testing.T) {
		id := byte(99) // Wrong ID
		payload := []byte{0x01, 0x80}

		report, err := NewBatteryReport(byte(id), payload)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "malformed battery stats report id")
		assert.Contains(t, err.Error(), "99")
		assert.Equal(t, BatteryReport{}, report)
	})

	t.Run("error wrong report id with different value", func(t *testing.T) {
		id := ReportKeyPress // Wrong ID (keypress instead of battery)
		payload := []byte{0x01, 0x80}

		report, err := NewBatteryReport(byte(id), payload)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "malformed battery stats report id")
		assert.Equal(t, BatteryReport{}, report)
	})
}

// TestNewKeyPressReport tests the NewKeyPressReport function
func TestNewKeyPressReport(t *testing.T) {
	t.Run("success single key", func(t *testing.T) {
		id := ReportKeyPress
		// Get a valid key ID from the keys package
		keysByName := keys.ByName()
		cam1Key := keysByName[keys.CAM1]
		payload := make([]byte, 2)
		binary.LittleEndian.PutUint16(payload, cam1Key.Id)

		report, err := NewKeyPressReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		require.Len(t, report.Keys, 1)
		assert.Equal(t, keys.CAM1, report.Keys[0].Name)
	})

	t.Run("success multiple keys", func(t *testing.T) {
		id := ReportKeyPress
		keysByName := keys.ByName()
		
		// Create payload with multiple keys
		payload := make([]byte, 6) // 3 keys * 2 bytes each
		binary.LittleEndian.PutUint16(payload[0:2], keysByName[keys.CAM1].Id)
		binary.LittleEndian.PutUint16(payload[2:4], keysByName[keys.CAM2].Id)
		binary.LittleEndian.PutUint16(payload[4:6], keysByName[keys.CAM3].Id)

		report, err := NewKeyPressReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		require.Len(t, report.Keys, 3)
		assert.Equal(t, keys.CAM1, report.Keys[0].Name)
		assert.Equal(t, keys.CAM2, report.Keys[1].Name)
		assert.Equal(t, keys.CAM3, report.Keys[2].Name)
	})

	t.Run("success no keys", func(t *testing.T) {
		id := ReportKeyPress
		payload := []byte{} // Empty payload

		report, err := NewKeyPressReport(byte(id), payload)

		require.NoError(t, err)
		assert.Equal(t, byte(id), report.Id)
		assert.Empty(t, report.Keys)
	})

	t.Run("success with jog keys", func(t *testing.T) {
		id := ReportKeyPress
		keysByName := keys.ByName()
		
		// Test with jog mode keys
		payload := make([]byte, 6)
		binary.LittleEndian.PutUint16(payload[0:2], keysByName[keys.JOG].Id)
		binary.LittleEndian.PutUint16(payload[2:4], keysByName[keys.SHTL].Id)
		binary.LittleEndian.PutUint16(payload[4:6], keysByName[keys.SCRL].Id)

		report, err := NewKeyPressReport(byte(id), payload)

		require.NoError(t, err)
		require.Len(t, report.Keys, 3)
		assert.Contains(t, report.Keys, keysByName[keys.JOG])
		assert.Contains(t, report.Keys, keysByName[keys.SHTL])
		assert.Contains(t, report.Keys, keysByName[keys.SCRL])
	})

	t.Run("success ignores unknown key ids", func(t *testing.T) {
		id := ReportKeyPress
		keysByName := keys.ByName()
		
		// Mix of valid and invalid key IDs
		payload := make([]byte, 4)
		binary.LittleEndian.PutUint16(payload[0:2], keysByName[keys.CAM1].Id)
		binary.LittleEndian.PutUint16(payload[2:4], uint16(9999)) // Invalid ID

		report, err := NewKeyPressReport(byte(id), payload)

		require.NoError(t, err)
		require.Len(t, report.Keys, 1) // Only the valid key should be included
		assert.Equal(t, keys.CAM1, report.Keys[0].Name)
	})

	t.Run("success handles even length payload with trailing zeros", func(t *testing.T) {
		id := ReportKeyPress
		keysByName := keys.ByName()

		// Even length payload with a zero key ID at the end
		payload := make([]byte, 4)
		binary.LittleEndian.PutUint16(payload[0:2], keysByName[keys.CAM1].Id)
		binary.LittleEndian.PutUint16(payload[2:4], uint16(0)) // Zero/invalid ID

		report, err := NewKeyPressReport(byte(id), payload)

		require.NoError(t, err)
		require.Len(t, report.Keys, 1) // Only the valid key should be included
		assert.Equal(t, keys.CAM1, report.Keys[0].Name)
	})

	t.Run("error wrong report id", func(t *testing.T) {
		id := byte(99) // Wrong ID
		payload := []byte{0x01, 0x00}

		report, err := NewKeyPressReport(byte(id), payload)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "malformed keypress input report id")
		assert.Contains(t, err.Error(), "99")
		assert.Equal(t, KeyPressReport{}, report)
	})

	t.Run("error wrong report id battery", func(t *testing.T) {
		id := ReportBattery // Wrong ID (battery instead of keypress)
		payload := []byte{0x01, 0x00}

		report, err := NewKeyPressReport(byte(id), payload)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "malformed keypress input report id")
		assert.Equal(t, KeyPressReport{}, report)
	})
}

// TestNewJogReport tests the NewJogReport function
func TestNewJogReport(t *testing.T) {
	t.Run("success relative mode", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x00 // RELATIVE mode
		binary.LittleEndian.PutUint32(payload[1:5], uint32(100))
		payload[5] = 0x00

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, "RELATIVE", report.Mode.Name)
		assert.Equal(t, int32(100), report.Value)
		assert.Equal(t, uint8(0), report.Unknown)
	})

	t.Run("success absolute mode", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x01 // ABSOLUTE mode
		binary.LittleEndian.PutUint32(payload[1:5], uint32(2048))
		payload[5] = 0xFF

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, "ABSOLUTE", report.Mode.Name)
		assert.Equal(t, int32(2048), report.Value)
		assert.Equal(t, uint8(0xFF), report.Unknown)
	})

	t.Run("success relative_2 mode", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x02 // RELATIVE_2 mode
		binary.LittleEndian.PutUint32(payload[1:5], uint32(500))

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, "RELATIVE_2", report.Mode.Name)
		assert.Equal(t, int32(500), report.Value)
	})

	t.Run("success absolute_deadzone mode", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x03 // ABSOLUTE_DEADZONE mode
		binary.LittleEndian.PutUint32(payload[1:5], uint32(4096))

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, "ABSOLUTE_DEADZONE", report.Mode.Name)
		assert.Equal(t, int32(4096), report.Value)
	})

	t.Run("success negative value", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x01 // ABSOLUTE mode
		binary.LittleEndian.PutUint32(payload[1:5], uint32(0xFFFFF000)) // Negative value

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Less(t, report.Value, int32(0))
	})

	t.Run("success zero value", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x00
		binary.LittleEndian.PutUint32(payload[1:5], uint32(0))

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, int32(0), report.Value)
	})

	t.Run("success max value", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x00
		binary.LittleEndian.PutUint32(payload[1:5], uint32(0xFFFFFFFF))

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, int32(-1), report.Value)
	})

	t.Run("unknown mode id returns empty mode", func(t *testing.T) {
		id := ReportJog
		payload := make([]byte, 6)
		payload[0] = 0x99 // Unknown mode ID
		binary.LittleEndian.PutUint32(payload[1:5], uint32(100))

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, "", report.Mode.Name)
		assert.Equal(t, 0, report.Mode.Id)
		assert.Equal(t, int32(100), report.Value)
	})

	t.Run("prints error for wrong report id but still returns report", func(t *testing.T) {
		id := byte(99) // Wrong ID
		payload := make([]byte, 6)
		payload[0] = 0x00
		binary.LittleEndian.PutUint32(payload[1:5], uint32(100))

		// Note: NewJogReport doesn't return an error, it just prints
		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, int32(100), report.Value)
	})

	t.Run("prints error for battery report id", func(t *testing.T) {
		id := ReportBattery // Wrong ID
		payload := make([]byte, 6)
		payload[0] = 0x00
		binary.LittleEndian.PutUint32(payload[1:5], uint32(100))

		report := NewJogReport(byte(id), payload)

		assert.Equal(t, byte(id), report.Id)
		assert.Equal(t, int32(100), report.Value)
	})
}

// TestToReport tests the ReportBytes.ToReport method
func TestToReport(t *testing.T) {
	t.Run("jog report", func(t *testing.T) {
		data := make([]byte, 7)
		data[0] = ReportJog
		data[1] = 0x00 // RELATIVE mode
		binary.LittleEndian.PutUint32(data[2:6], uint32(100))

		report, err := ReportBytes(data).ToReport()

		require.NoError(t, err)
		require.NotNil(t, report)
		jogReport, ok := report.(JogReport)
		require.True(t, ok)
		assert.Equal(t, byte(ReportJog), jogReport.Id)
		assert.Equal(t, int32(100), jogReport.Value)
	})

	t.Run("keypress report", func(t *testing.T) {
		keysByName := keys.ByName()
		data := make([]byte, 3)
		data[0] = ReportKeyPress
		binary.LittleEndian.PutUint16(data[1:3], keysByName[keys.CAM1].Id)

		report, err := ReportBytes(data).ToReport()

		require.NoError(t, err)
		require.NotNil(t, report)
		keyReport, ok := report.(KeyPressReport)
		require.True(t, ok)
		assert.Equal(t, byte(ReportKeyPress), keyReport.Id)
		require.Len(t, keyReport.Keys, 1)
		assert.Equal(t, keys.CAM1, keyReport.Keys[0].Name)
	})

	t.Run("battery report", func(t *testing.T) {
		data := []byte{ReportBattery, 0x01, 0x80}

		report, err := ReportBytes(data).ToReport()

		require.NoError(t, err)
		require.NotNil(t, report)
		batteryReport, ok := report.(BatteryReport)
		require.True(t, ok)
		assert.Equal(t, byte(ReportBattery), batteryReport.Id)
		assert.True(t, batteryReport.Charging)
		assert.InDelta(t, 0.5, batteryReport.Battery, 0.01)
	})

	t.Run("unknown report id returns error", func(t *testing.T) {
		data := []byte{99, 0x00, 0x00}

		report, err := ReportBytes(data).ToReport()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown report id")
		assert.Contains(t, err.Error(), "99")
		assert.Nil(t, report)
	})

	t.Run("empty data panics", func(t *testing.T) {
		data := []byte{}

		// ToReport doesn't handle empty data gracefully - it panics
		assert.Panics(t, func() {
			_, _ = ReportBytes(data).ToReport()
		})
	})

	t.Run("single byte panics", func(t *testing.T) {
		data := []byte{ReportJog}

		// ToReport with single byte causes panic in NewJogReport due to slice access
		assert.Panics(t, func() {
			_, _ = ReportBytes(data).ToReport()
		})
	})
}

// TestReportConstants tests the report ID constants
func TestReportConstants(t *testing.T) {
	t.Run("report jog constant", func(t *testing.T) {
		assert.Equal(t, byte(3), byte(ReportJog))
	})

	t.Run("report keypress constant", func(t *testing.T) {
		assert.Equal(t, byte(4), byte(ReportKeyPress))
	})

	t.Run("report battery constant", func(t *testing.T) {
		assert.Equal(t, byte(7), byte(ReportBattery))
	})

	t.Run("report constants are unique", func(t *testing.T) {
		assert.NotEqual(t, ReportJog, ReportKeyPress)
		assert.NotEqual(t, ReportJog, ReportBattery)
		assert.NotEqual(t, ReportKeyPress, ReportBattery)
	})
}