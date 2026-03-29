package speedEditor

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"github.com/JamesBalazs/speed-editor-client/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupAuthFixture creates an AuthHandler with a mocked HID device
func setupAuthFixture(t *testing.T) (*AuthHandler, *MockHIDDevice) {
	mockDevice := new(MockHIDDevice)

	ah := AuthHandler{
		device: mockDevice,
	}

	return &ah, mockDevice
}

// TestResetAuthState tests the ResetAuthState method
func TestResetAuthState(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()

		err := ah.ResetAuthState()

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error sending feature report", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(0, errors.New("send error")).Once()

		err := ah.ResetAuthState()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send feature report")
		assert.Contains(t, err.Error(), "send error")
		mockDevice.AssertExpectations(t)
	})
}

// TestGetKeyboardChallenge tests the GetKeyboardChallenge method
func TestGetKeyboardChallenge(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)
		expectedChallenge := uint64(1234567890)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], expectedChallenge)
		}).Return(10, nil).Once()

		challenge, err := ah.GetKeyboardChallenge()

		require.NoError(t, err)
		assert.Equal(t, expectedChallenge, challenge)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error getting feature report", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("get error")).Once()

		challenge, err := ah.GetKeyboardChallenge()

		require.Error(t, err)
		assert.Equal(t, uint64(0), challenge)
		assert.Contains(t, err.Error(), "failed to get feature report")
		assert.Contains(t, err.Error(), "get error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("unexpected header returns error", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			// Wrong header
			buf[0] = 0x06
			buf[1] = 0x99
		}).Return(10, nil).Once()

		challenge, err := ah.GetKeyboardChallenge()

		require.Error(t, err)
		assert.Equal(t, uint64(0), challenge)
		assert.Contains(t, err.Error(), "unexpected keyboard challenge response header")
		mockDevice.AssertExpectations(t)
	})

	t.Run("unexpected header with different wrong values returns error", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			// Wrong header - different values
			buf[0] = 0x00
			buf[1] = 0x00
		}).Return(10, nil).Once()

		challenge, err := ah.GetKeyboardChallenge()

		require.Error(t, err)
		assert.Equal(t, uint64(0), challenge)
		assert.Contains(t, err.Error(), "unexpected keyboard challenge response header")
		mockDevice.AssertExpectations(t)
	})
}

// TestSendHostChallenge tests the SendHostChallenge method
func TestSendHostChallenge(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()

		err := ah.SendHostChallenge()

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error sending feature report", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(0, errors.New("send error")).Once()

		err := ah.SendHostChallenge()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send feature report")
		assert.Contains(t, err.Error(), "send error")
		mockDevice.AssertExpectations(t)
	})
}

// TestGetHostChallengeResponse tests the GetHostChallengeResponse method
func TestGetHostChallengeResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x02
			buf[2] = 0x01
			buf[3] = 0x02
		}).Return(10, nil).Once()

		response, err := ah.GetHostChallengeResponse()

		require.NoError(t, err)
		assert.Len(t, response, 10)
		assert.Equal(t, byte(0x06), response[0])
		assert.Equal(t, byte(0x02), response[1])
		mockDevice.AssertExpectations(t)
	})

	t.Run("error getting feature report", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("get error")).Once()

		response, err := ah.GetHostChallengeResponse()

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get feature report")
		assert.Contains(t, err.Error(), "get error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("unexpected header returns error", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			// Wrong header
			buf[0] = 0x06
			buf[1] = 0x99
		}).Return(10, nil).Once()

		response, err := ah.GetHostChallengeResponse()

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "unexpected host challenge response header")
		mockDevice.AssertExpectations(t)
	})
}

// TestSendAuthChallengeResponse tests the SendAuthChallengeResponse method
func TestSendAuthChallengeResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)
		expectedResponse := uint64(9876543210)

		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			// Verify the response bytes are correct
			assert.Equal(t, byte(0x06), buf[0])
			assert.Equal(t, byte(0x03), buf[1])
			actualResponse := binary.LittleEndian.Uint64(buf[2:])
			assert.Equal(t, expectedResponse, actualResponse)
		}).Return(10, nil).Once()

		err := ah.SendAuthChallengeResponse(expectedResponse)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success with zero response", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil).Once()

		err := ah.SendAuthChallengeResponse(0)

		require.NoError(t, err)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error sending feature report", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("send error")).Once()

		err := ah.SendAuthChallengeResponse(12345)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send feature report")
		assert.Contains(t, err.Error(), "send error")
		mockDevice.AssertExpectations(t)
	})
}

// TestGetAuthChallengeResult tests the GetAuthChallengeResult method
func TestGetAuthChallengeResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)
		expectedResult := uint16(65535)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x04
			binary.LittleEndian.PutUint16(buf[2:4], expectedResult)
		}).Return(10, nil).Once()

		result, err := ah.GetAuthChallengeResult()

		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success with different timeout value", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)
		expectedResult := uint16(3600)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x04
			binary.LittleEndian.PutUint16(buf[2:4], expectedResult)
		}).Return(10, nil).Once()

		result, err := ah.GetAuthChallengeResult()

		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error getting feature report", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("get error")).Once()

		result, err := ah.GetAuthChallengeResult()

		require.Error(t, err)
		assert.Equal(t, uint16(0), result)
		assert.Contains(t, err.Error(), "failed to get feature report")
		assert.Contains(t, err.Error(), "get error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("unexpected header returns error", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			// Wrong header
			buf[0] = 0x06
			buf[1] = 0x99
		}).Return(10, nil).Once()

		result, err := ah.GetAuthChallengeResult()

		require.Error(t, err)
		assert.Equal(t, uint16(0), result)
		assert.Contains(t, err.Error(), "unexpected auth response header")
		mockDevice.AssertExpectations(t)
	})
}

// TestAuthHandlerAuthenticate tests the full Authenticate method
func TestAuthHandlerAuthenticate(t *testing.T) {
	t.Run("success full flow", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)
		keyboardChallenge := uint64(1234567890)
		expectedReauthTimeout := uint16(65535)

		// Step 1: ResetAuthState - SendFeatureReport with default state
		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()

		// Step 2: GetKeyboardChallenge - GetFeatureReport with keyboard challenge header
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], keyboardChallenge)
		}).Return(10, nil).Once()

		// Step 3: SendHostChallenge - SendFeatureReport with host challenge
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()

		// Step 4: GetHostChallengeResponse - GetFeatureReport with host challenge response header
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x02
		}).Return(10, nil).Once()

		// Step 5: SendAuthChallengeResponse - SendFeatureReport with calculated response
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			// Verify header
			assert.Equal(t, byte(0x06), buf[0])
			assert.Equal(t, byte(0x03), buf[1])
		}).Return(10, nil).Once()

		// Step 6: GetAuthChallengeResult - GetFeatureReport with auth response header
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x04
			binary.LittleEndian.PutUint16(buf[2:4], expectedReauthTimeout)
		}).Return(10, nil).Once()

		reauthDuration, err := ah.Authenticate()

		require.NoError(t, err)
		// Verify the reauth duration is calculated correctly (timeout - 10 seconds)
		expectedDuration := time.Duration(expectedReauthTimeout-10) * time.Second
		assert.Equal(t, expectedDuration, reauthDuration)
		mockDevice.AssertExpectations(t)
	})

	t.Run("success with different challenge values", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)
		keyboardChallenge := uint64(9999999999)
		expectedReauthTimeout := uint16(3600)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], keyboardChallenge)
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x02
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x04
			binary.LittleEndian.PutUint16(buf[2:4], expectedReauthTimeout)
		}).Return(10, nil).Once()

		reauthDuration, err := ah.Authenticate()

		require.NoError(t, err)
		expectedDuration := time.Duration(expectedReauthTimeout-10) * time.Second
		assert.Equal(t, expectedDuration, reauthDuration)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on ResetAuthState", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(0, errors.New("send error")).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to reset auth state")
		assert.Contains(t, err.Error(), "send error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on GetKeyboardChallenge", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("get error")).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to get keyboard challenge")
		assert.Contains(t, err.Error(), "get error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on GetKeyboardChallenge unexpected header", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x99 // Wrong header
		}).Return(10, nil).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to get keyboard challenge")
		assert.Contains(t, err.Error(), "unexpected keyboard challenge response header")
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on SendHostChallenge", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(0, errors.New("send error")).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to send host challenge")
		assert.Contains(t, err.Error(), "send error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on GetHostChallengeResponse", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("get error")).Once()
		// GetHostChallengeResponse error is ignored, so Authenticate continues to SendAuthChallengeResponse
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x04
			binary.LittleEndian.PutUint16(buf[2:4], uint16(65535))
		}).Return(10, nil).Once()

		reauthDuration, err := ah.Authenticate()

		// GetHostChallengeResponse error is ignored in Authenticate
		require.NoError(t, err)
		assert.NotEqual(t, time.Duration(0), reauthDuration)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on GetHostChallengeResponse unexpected header", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x99 // Wrong header
		}).Return(10, nil).Once()
		// GetHostChallengeResponse error is ignored, so Authenticate continues to SendAuthChallengeResponse
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x04
			binary.LittleEndian.PutUint16(buf[2:4], uint16(65535))
		}).Return(10, nil).Once()

		reauthDuration, err := ah.Authenticate()

		// GetHostChallengeResponse error is ignored in Authenticate
		require.NoError(t, err)
		assert.NotEqual(t, time.Duration(0), reauthDuration)
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on SendAuthChallengeResponse", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x02
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("send error")).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to send auth challenge response")
		assert.Contains(t, err.Error(), "send error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on GetAuthChallengeResult", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x02
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Return(0, errors.New("get error")).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to get auth challenge result")
		assert.Contains(t, err.Error(), "get error")
		mockDevice.AssertExpectations(t)
	})

	t.Run("error on GetAuthChallengeResult unexpected header", func(t *testing.T) {
		ah, mockDevice := setupAuthFixture(t)

		mockDevice.On("SendFeatureReport", featureReportDefaultState).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x00
			binary.LittleEndian.PutUint64(buf[2:], uint64(12345))
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", featureReportHostChallenge).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x02
		}).Return(10, nil).Once()
		mockDevice.On("SendFeatureReport", mock.AnythingOfType("[]uint8")).Return(10, nil).Once()
		mockDevice.On("GetFeatureReport", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
			buf := args.Get(0).([]byte)
			buf[0] = 0x06
			buf[1] = 0x99 // Wrong header
		}).Return(10, nil).Once()

		reauthDuration, err := ah.Authenticate()

		require.Error(t, err)
		assert.Equal(t, time.Duration(0), reauthDuration)
		assert.Contains(t, err.Error(), "failed to get auth challenge result")
		assert.Contains(t, err.Error(), "unexpected auth response header")
		mockDevice.AssertExpectations(t)
	})
}

// TestAuthChallengeCalculation verifies the challenge response calculation is correct
func TestAuthChallengeCalculation(t *testing.T) {
	t.Run("calculate challenge response", func(t *testing.T) {
		// Test that the auth.CalculateChallengeResponse function works
		challenge := uint64(1234567890)
		response := auth.CalculateChallengeResponse(challenge)

		// The response should be different from the challenge
		assert.NotEqual(t, challenge, response)
		// The response should be non-zero
		assert.NotZero(t, response)
	})

	t.Run("calculate challenge response with zero", func(t *testing.T) {
		challenge := uint64(0)
		response := auth.CalculateChallengeResponse(challenge)

		// Response should be non-zero even with zero challenge
		assert.NotZero(t, response)
	})

	t.Run("calculate challenge response with max value", func(t *testing.T) {
		challenge := uint64(18446744073709551615) // Max uint64
		response := auth.CalculateChallengeResponse(challenge)

		// Response should be non-zero
		assert.NotZero(t, response)
	})
}

// TestAuthHandlerInterface verifies AuthHandler implements AuthHandlerInt
func TestAuthHandlerInterface(t *testing.T) {
	t.Run("AuthHandler implements AuthHandlerInt", func(t *testing.T) {
		var _ AuthHandlerInt = (*AuthHandler)(nil)
	})
}
