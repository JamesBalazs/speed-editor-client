package speedEditor

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"github.com/JamesBalazs/speed-editor-client/input"
	"github.com/JamesBalazs/speed-editor-client/keys"
	"github.com/sstallion/go-hid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockHIDDevice is a mock implementation of deviceInterface for testing
type MockHIDDevice struct {
	mock.Mock
}

func (m *MockHIDDevice) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockHIDDevice) Read(buf []byte) (int, error) {
	args := m.Called(buf)
	return args.Int(0), args.Error(1)
}

func (m *MockHIDDevice) Write(buf []byte) (int, error) {
	args := m.Called(buf)
	return args.Int(0), args.Error(1)
}

func (m *MockHIDDevice) GetDeviceInfo() (*hid.DeviceInfo, error) {
	args := m.Called()
	return args.Get(0).(*hid.DeviceInfo), args.Error(1)
}

func (m *MockHIDDevice) GetFeatureReport(buf []byte) (int, error) {
	args := m.Called(buf)
	return args.Int(0), args.Error(1)
}

func (m *MockHIDDevice) SendFeatureReport(buf []byte) (int, error) {
	args := m.Called(buf)
	return args.Int(0), args.Error(1)
}

// setupFixture creates a SpeedEditor instance with a mocked HID device
func setupFixture(t *testing.T) (*SpeedEditor, *MockHIDDevice) {
	mockDevice := new(MockHIDDevice)
	deviceInfo := &hid.DeviceInfo{
		MfrStr:     "Test Manufacturer",
		ProductStr: "Test Product",
		SerialNbr:  "TEST123",
	}

	se := &SpeedEditor{
		device:     mockDevice,
		deviceInfo: *deviceInfo,
	}

	return se, mockDevice
}

// setupFixtureWithDeviceInfo creates a SpeedEditor with custom device info
func setupFixtureWithDeviceInfo(t *testing.T, deviceInfo *hid.DeviceInfo) (*SpeedEditor, *MockHIDDevice) {
	mockDevice := new(MockHIDDevice)

	se := &SpeedEditor{
		device:     mockDevice,
		deviceInfo: *deviceInfo,
	}

	return se, mockDevice
}

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	t.Run("connects to device or returns error", func(t *testing.T) {
		// This test will succeed if a physical device is connected,
		// or fail with an error if no device is present
		client, err := NewClient()
		
		if err != nil {
			// No device connected - verify error message
			assert.Contains(t, err.Error(), "failed to create client")
			assert.Nil(t, client)
		} else {
			// Device connected - verify client is valid
			assert.NotNil(t, client)
			deviceInfo := client.GetDeviceInfo()
			assert.NotEmpty(t, deviceInfo.SerialNbr)
		}
	})
}

// TestInitialize tests the initialize method
func TestInitialize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDevice := new(MockHIDDevice)
		deviceInfo := &hid.DeviceInfo{
			MfrStr:     "Test Manufacturer",
			ProductStr: "Test Product",
			SerialNbr:  "TEST123",
		}
		mockDevice.On("GetDeviceInfo").Return(deviceInfo, nil).Once()

		se := &SpeedEditor{
			device:      mockDevice,
			AuthHandler: AuthHandler{},
		}

		err := se.initialize()

		require.NoError(t, err)
		assert.NotNil(t, se.jogHandler)
		assert.NotNil(t, se.batteryHandler)
		assert.NotNil(t, se.keyPressHandler)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error getting device info", func(t *testing.T) {
		mockDevice := new(MockHIDDevice)
		mockDevice.On("GetDeviceInfo").Return((*hid.DeviceInfo)(nil), errors.New("device info error")).Once()

		se := &SpeedEditor{
			device:      mockDevice,
			AuthHandler: AuthHandler{},
		}

		err := se.initialize()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get device info")
		assert.Contains(t, err.Error(), "device info error")
		mockDevice.AssertExpectations(t)
	})
}

// TestGetDeviceInfo tests the GetDeviceInfo method
func TestGetDeviceInfo(t *testing.T) {
	t.Run("success returns cached device info", func(t *testing.T) {
		deviceInfo := &hid.DeviceInfo{
			MfrStr:    "Test Manufacturer",
			ProductStr: "Test Product",
			SerialNbr: "TEST123",
		}

		se := &SpeedEditor{
			deviceInfo: *deviceInfo,
		}

		result := se.GetDeviceInfo()
		assert.Equal(t, *deviceInfo, result)
	})
}

// TestRead tests the Read method
func TestRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		se, mockDevice := setupFixture(t)
		expectedData := []byte{0x04, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00}

		mockDevice.On("Read", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			copy(buf, expectedData)
		}).Return(len(expectedData), nil).Once()

		data, n, err := se.Read()

		require.NoError(t, err)
		assert.Equal(t, len(expectedData), n)
		assert.Equal(t, expectedData, data)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error reading from device", func(t *testing.T) {
		se, mockDevice := setupFixture(t)
		mockDevice.On("Read", mock.AnythingOfType("[]uint8")).Return(0, errors.New("read error")).Once()

		data, n, err := se.Read()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read from device")
		assert.Contains(t, err.Error(), "read error")
		assert.Nil(t, data)
		assert.Equal(t, 0, n)
		mockDevice.AssertExpectations(t)
	})
}

// TestSetLeds tests the SetLeds method
func TestSetLeds(t *testing.T) {
	t.Run("success single led", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(LedReportId), buf[0])
			expectedBitmask := uint32(0x00000001)
			actualBitmask := binary.LittleEndian.Uint32(buf[1:5])
			assert.Equal(t, expectedBitmask, actualBitmask)
		}).Return(5, nil).Once()

		leds := []uint32{0x01}
		err := se.SetLeds(leds)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success multiple leds", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(LedReportId), buf[0])
			// Verify OR operation: 0x01 | 0x02 | 0x04 = 0x07
			expectedBitmask := uint32(0x07)
			actualBitmask := binary.LittleEndian.Uint32(buf[1:5])
			assert.Equal(t, expectedBitmask, actualBitmask)
		}).Return(5, nil).Once()

		leds := []uint32{0x01, 0x02, 0x04}
		err := se.SetLeds(leds)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success empty leds", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(LedReportId), buf[0])
			expectedBitmask := uint32(0x00000000)
			actualBitmask := binary.LittleEndian.Uint32(buf[1:5])
			assert.Equal(t, expectedBitmask, actualBitmask)
		}).Return(5, nil).Once()

		leds := []uint32{}
		err := se.SetLeds(leds)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error writing to device", func(t *testing.T) {
		se, mockDevice := setupFixture(t)
		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Return(0, errors.New("write error")).Once()

		leds := []uint32{0x01}
		err := se.SetLeds(leds)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set LEDs")
		assert.Contains(t, err.Error(), "write error")
		mockDevice.AssertExpectations(t)
	})
}

// TestSetJogMode tests the SetJogMode method
func TestSetJogMode(t *testing.T) {
	t.Run("success absolute mode", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(JogModeReportId), buf[0])
			assert.Equal(t, byte(1), buf[1]) // ABSOLUTE mode
			assert.Equal(t, byte(0), buf[2])
			assert.Equal(t, byte(255), buf[6])
		}).Return(7, nil).Once()

		err := se.SetJogMode(1)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success relative mode", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(JogModeReportId), buf[0])
			assert.Equal(t, byte(0), buf[1]) // RELATIVE mode
		}).Return(7, nil).Once()

		err := se.SetJogMode(0)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error writing to device", func(t *testing.T) {
		se, mockDevice := setupFixture(t)
		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Return(0, errors.New("write error")).Once()

		err := se.SetJogMode(1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set jog mode")
		assert.Contains(t, err.Error(), "write error")
		mockDevice.AssertExpectations(t)
	})
}

// TestSetJogLeds tests the SetJogLeds method
func TestSetJogLeds(t *testing.T) {
	t.Run("success single jog led", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(JogLedReportId), buf[0])
			assert.Equal(t, byte(0x01), buf[1])
		}).Return(2, nil).Once()

		leds := []uint8{0x01}
		err := se.SetJogLeds(leds)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success multiple jog leds", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(JogLedReportId), buf[0])
			// Verify OR operation: 0x01 | 0x02 | 0x04 = 0x07
			expectedBitmask := uint8(0x07)
			actualBitmask := buf[1]
			assert.Equal(t, expectedBitmask, actualBitmask)
		}).Return(2, nil).Once()

		leds := []uint8{0x01, 0x02, 0x04}
		err := se.SetJogLeds(leds)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success empty jog leds", func(t *testing.T) {
		se, mockDevice := setupFixture(t)

		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			assert.Equal(t, byte(JogLedReportId), buf[0])
			assert.Equal(t, byte(0x00), buf[1])
		}).Return(2, nil).Once()

		leds := []uint8{}
		err := se.SetJogLeds(leds)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error writing to device", func(t *testing.T) {
		se, mockDevice := setupFixture(t)
		mockDevice.On("Write", mock.AnythingOfType("[]uint8")).Return(0, errors.New("write error")).Once()

		leds := []uint8{0x01}
		err := se.SetJogLeds(leds)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set jog LEDs")
		assert.Contains(t, err.Error(), "write error")
		mockDevice.AssertExpectations(t)
	})
}

// TestSetJogHandler tests the SetJogHandler method
func TestSetJogHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		se := &SpeedEditor{}
		called := false
		handler := func(client SpeedEditorInt, report input.JogReport) {
			called = true
		}

		se.SetJogHandler(handler)
		assert.NotNil(t, se.jogHandler)

		se.jogHandler(se, input.JogReport{})
		assert.True(t, called)
	})

	t.Run("success replace handler", func(t *testing.T) {
		se := &SpeedEditor{}
		firstCalled := false
		secondCalled := false

		firstHandler := func(client SpeedEditorInt, report input.JogReport) {
			firstCalled = true
		}
		secondHandler := func(client SpeedEditorInt, report input.JogReport) {
			secondCalled = true
		}

		se.SetJogHandler(firstHandler)
		se.SetJogHandler(secondHandler)

		se.jogHandler(se, input.JogReport{})
		assert.False(t, firstCalled)
		assert.True(t, secondCalled)
	})
}

// TestSetBatteryHandler tests the SetBatteryHandler method
func TestSetBatteryHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		se := &SpeedEditor{}
		called := false
		handler := func(client SpeedEditorInt, report input.BatteryReport) {
			called = true
		}

		se.SetBatteryHandler(handler)
		assert.NotNil(t, se.batteryHandler)

		se.batteryHandler(se, input.BatteryReport{})
		assert.True(t, called)
	})

	t.Run("success replace handler", func(t *testing.T) {
		se := &SpeedEditor{}
		firstCalled := false
		secondCalled := false

		firstHandler := func(client SpeedEditorInt, report input.BatteryReport) {
			firstCalled = true
		}
		secondHandler := func(client SpeedEditorInt, report input.BatteryReport) {
			secondCalled = true
		}

		se.SetBatteryHandler(firstHandler)
		se.SetBatteryHandler(secondHandler)

		se.batteryHandler(se, input.BatteryReport{})
		assert.False(t, firstCalled)
		assert.True(t, secondCalled)
	})
}

// TestSetKeyPressHandler tests the SetKeyPressHandler method
func TestSetKeyPressHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		se := &SpeedEditor{}
		called := false
		handler := func(client SpeedEditorInt, report input.KeyPressReport) {
			called = true
		}

		se.SetKeyPressHandler(handler)
		assert.NotNil(t, se.keyPressHandler)

		se.keyPressHandler(se, input.KeyPressReport{})
		assert.True(t, called)
	})

	t.Run("success replace handler", func(t *testing.T) {
		se := &SpeedEditor{}
		firstCalled := false
		secondCalled := false

		firstHandler := func(client SpeedEditorInt, report input.KeyPressReport) {
			firstCalled = true
		}
		secondHandler := func(client SpeedEditorInt, report input.KeyPressReport) {
			secondCalled = true
		}

		se.SetKeyPressHandler(firstHandler)
		se.SetKeyPressHandler(secondHandler)

		se.keyPressHandler(se, input.KeyPressReport{})
		assert.False(t, firstCalled)
		assert.True(t, secondCalled)
	})
}

// TestHandleReport tests the HandleReport method
func TestHandleReport(t *testing.T) {
	t.Run("jog report", func(t *testing.T) {
		se := &SpeedEditor{}
		jogCalled := false
		se.SetJogHandler(func(client SpeedEditorInt, report input.JogReport) {
			jogCalled = true
			assert.Equal(t, int32(100), report.Value)
		})

		report := input.JogReport{Value: 100}
		se.HandleReport(report)
		assert.True(t, jogCalled)
	})

	t.Run("battery report", func(t *testing.T) {
		se := &SpeedEditor{}
		batteryCalled := false
		se.SetBatteryHandler(func(client SpeedEditorInt, report input.BatteryReport) {
			batteryCalled = true
			assert.Equal(t, float32(0.5), report.Battery)
		})

		report := input.BatteryReport{Battery: 0.5}
		se.HandleReport(report)
		assert.True(t, batteryCalled)
	})

	t.Run("keypress report", func(t *testing.T) {
		se := &SpeedEditor{}
		keyPressCalled := false
		se.SetKeyPressHandler(func(client SpeedEditorInt, report input.KeyPressReport) {
			keyPressCalled = true
			assert.Len(t, report.Keys, 1)
		})

		keysByName := keys.ByName()
		report := input.KeyPressReport{Keys: []keys.Key{keysByName[keys.CAM1]}}
		se.HandleReport(report)
		assert.True(t, keyPressCalled)
	})
}

// TestHandleJog tests the HandleJog method
func TestHandleJog(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		se := &SpeedEditor{}
		handlerCalled := false
		var receivedReport input.JogReport

		se.SetJogHandler(func(client SpeedEditorInt, report input.JogReport) {
			handlerCalled = true
			receivedReport = report
		})

		expectedReport := input.JogReport{
			Id:      3,
			Value:   500,
			Unknown: 0,
		}

		se.HandleJog(expectedReport)
		assert.True(t, handlerCalled)
		assert.Equal(t, expectedReport, receivedReport)
	})
}

// TestHandleBattery tests the HandleBattery method
func TestHandleBattery(t *testing.T) {
	t.Run("success charging", func(t *testing.T) {
		se := &SpeedEditor{}
		handlerCalled := false
		var receivedReport input.BatteryReport

		se.SetBatteryHandler(func(client SpeedEditorInt, report input.BatteryReport) {
			handlerCalled = true
			receivedReport = report
		})

		expectedReport := input.BatteryReport{
			Id:       7,
			Charging: true,
			Battery:  0.75,
		}

		se.HandleBattery(expectedReport)
		assert.True(t, handlerCalled)
		assert.Equal(t, expectedReport, receivedReport)
	})

	t.Run("success not charging", func(t *testing.T) {
		se := &SpeedEditor{}
		handlerCalled := false
		var receivedReport input.BatteryReport

		se.SetBatteryHandler(func(client SpeedEditorInt, report input.BatteryReport) {
			handlerCalled = true
			receivedReport = report
		})

		expectedReport := input.BatteryReport{
			Id:       7,
			Charging: false,
			Battery:  0.25,
		}

		se.HandleBattery(expectedReport)
		assert.True(t, handlerCalled)
		assert.Equal(t, expectedReport, receivedReport)
	})
}

// TestHandleKeyPress tests the HandleKeyPress method
func TestHandleKeyPress(t *testing.T) {
	t.Run("single key", func(t *testing.T) {
		se := &SpeedEditor{}
		handlerCalled := false
		var receivedReport input.KeyPressReport

		se.SetKeyPressHandler(func(client SpeedEditorInt, report input.KeyPressReport) {
			handlerCalled = true
			receivedReport = report
		})

		keysByName := keys.ByName()
		expectedReport := input.KeyPressReport{
			Id:   4,
			Keys: []keys.Key{keysByName[keys.CAM1]},
		}

		se.HandleKeyPress(expectedReport)
		assert.True(t, handlerCalled)
		assert.Equal(t, expectedReport, receivedReport)
	})

	t.Run("multiple keys", func(t *testing.T) {
		se := &SpeedEditor{}
		handlerCalled := false
		var receivedReport input.KeyPressReport

		se.SetKeyPressHandler(func(client SpeedEditorInt, report input.KeyPressReport) {
			handlerCalled = true
			receivedReport = report
		})

		keysByName := keys.ByName()
		expectedReport := input.KeyPressReport{
			Id:   4,
			Keys: []keys.Key{keysByName[keys.CAM1], keysByName[keys.CAM2], keysByName[keys.CAM3]},
		}

		se.HandleKeyPress(expectedReport)
		assert.True(t, handlerCalled)
		assert.Equal(t, expectedReport, receivedReport)
	})

	t.Run("no keys", func(t *testing.T) {
		se := &SpeedEditor{}
		handlerCalled := false
		var receivedReport input.KeyPressReport

		se.SetKeyPressHandler(func(client SpeedEditorInt, report input.KeyPressReport) {
			handlerCalled = true
			receivedReport = report
		})

		expectedReport := input.KeyPressReport{
			Id:   4,
			Keys: []keys.Key{},
		}

		se.HandleKeyPress(expectedReport)
		assert.True(t, handlerCalled)
		assert.Equal(t, expectedReport, receivedReport)
	})
}

// TestAuthenticate tests the Authenticate method
func TestAuthenticate(t *testing.T) {
	t.Run("success starts goroutine", func(t *testing.T) {
		mockDevice := new(MockHIDDevice)

		// Setup mock expectations for the auth handshake
		// The auth flow has multiple GetFeatureReport calls with different expected headers
		callCount := 0
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil)
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			callCount++
			switch callCount {
			case 1:
				// Keyboard challenge response
				buf[0] = 0x06
				buf[1] = 0x00
				binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
			case 2:
				// Host challenge response
				buf[0] = 0x06
				buf[1] = 0x02
			case 3:
				// Auth challenge result
				buf[0] = 0x06
				buf[1] = 0x04
				binary.LittleEndian.PutUint16(buf[2:], uint16(65535))
			}
		}).Return(10, nil)

		se := &SpeedEditor{
			device:      mockDevice,
			AuthHandler: AuthHandler{device: mockDevice},
		}

		err := se.Authenticate()
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
	})
}

// TestDefaultHandlers tests the default handler functions
func TestDefaultHandlers(t *testing.T) {
	t.Run("default jog handler", func(t *testing.T) {
		se := &SpeedEditor{}
		report := input.JogReport{Value: 100}

		assert.NotPanics(t, func() {
			defaultJogHandler(se, report)
		})
	})

	t.Run("default battery handler", func(t *testing.T) {
		se := &SpeedEditor{}
		report := input.BatteryReport{
			Charging: true,
			Battery:  0.5,
		}

		assert.NotPanics(t, func() {
			defaultBatteryHandler(se, report)
		})
	})

	t.Run("default keypress handler", func(t *testing.T) {
		se := &SpeedEditor{}
		keysByName := keys.ByName()
		report := input.KeyPressReport{
			Keys: []keys.Key{keysByName[keys.CAM1]},
		}

		assert.NotPanics(t, func() {
			defaultKeyPressHandler(se, report)
		})
	})
}

// TestNullHandlers tests the null handler functions
func TestNullHandlers(t *testing.T) {
	t.Run("null jog handler", func(t *testing.T) {
		se := &SpeedEditor{}
		report := input.JogReport{Value: 100}

		assert.NotPanics(t, func() {
			NullJogHandler(se, report)
		})
	})

	t.Run("null battery handler", func(t *testing.T) {
		se := &SpeedEditor{}
		report := input.BatteryReport{Battery: 0.5}

		assert.NotPanics(t, func() {
			NullBatteryHandler(se, report)
		})
	})

	t.Run("null keypress handler", func(t *testing.T) {
		se := &SpeedEditor{}
		keysByName := keys.ByName()
		report := input.KeyPressReport{Keys: []keys.Key{keysByName[keys.CAM1]}}

		assert.NotPanics(t, func() {
			NullKeyPressHandler(se, report)
		})
	})
}
