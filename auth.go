// Based on the authentication algorithm by Sylvain Munaut <tnt@246tNt.com> (Apache-2.0)
// https://github.com/smunaut/blackmagic-misc/blob/master/bmd.py
// SPDX-License-Identifier: Apache-2.0

package speedEditor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/JamesBalazs/speed-editor-client/auth"
)

var (
	featureReportDefaultState  = []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	featureReportHostChallenge = []byte{0x06, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	authChallengeHeader = []byte{0x06, 0x03}

	expectedKeyboardChallengeResponseHeader = []byte{0x06, 0x00}
	expectedHostChallengeResponseHeader     = []byte{0x06, 0x02}
	expectedAuthResponseHeader              = []byte{0x06, 0x04}

	reauthTimeout = uint16(65535) // initialise to highest possible value
)

type AuthHandlerInt interface {
	Authenticate() (time.Duration, error)
	ResetAuthState() error
	GetKeyboardChallenge() (uint64, error)
	SendHostChallenge() error
	GetHostChallengeResponse() ([]byte, error)
	SendAuthChallengeResponse(response uint64) error
	GetAuthChallengeResult() (uint16, error)
}

type AuthHandler struct {
	device deviceInterface
}

// Authenticate handles the entire handshake between the host and the Speed Editor.
//
// It returns the duration before the Speed Editor expects a reauth.
func (ah AuthHandler) Authenticate() (time.Duration, error) {
	if err := ah.ResetAuthState(); err != nil {
		return 0, fmt.Errorf("failed to reset auth state: %w", err)
	}

	challenge, err := ah.GetKeyboardChallenge()
	if err != nil {
		return 0, fmt.Errorf("failed to get keyboard challenge: %w", err)
	}

	if err := ah.SendHostChallenge(); err != nil {
		return 0, fmt.Errorf("failed to send host challenge: %w", err)
	}

	// We don't care about the response or error, since we don't care if it's a real Speed Editor
	_, _ = ah.GetHostChallengeResponse()

	response := auth.CalculateChallengeResponse(challenge)
	if err := ah.SendAuthChallengeResponse(response); err != nil {
		return 0, fmt.Errorf("failed to send auth challenge response: %w", err)
	}

	reauthTimeout, err := ah.GetAuthChallengeResult()
	if err != nil {
		return 0, fmt.Errorf("failed to get auth challenge result: %w", err)
	}

	return time.Duration(reauthTimeout-10) * time.Second, nil
}

func (ah AuthHandler) ResetAuthState() error {
	_, err := ah.device.SendFeatureReport(featureReportDefaultState)
	if err != nil {
		return fmt.Errorf("failed to send feature report: %w", err)
	}
	return nil
}

func (ah AuthHandler) GetKeyboardChallenge() (uint64, error) {
	// get keyboard challenge, store in a new copy of the byte array
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := ah.device.GetFeatureReport(data)
	if err != nil {
		return 0, fmt.Errorf("failed to get feature report: %w", err)
	}

	if !bytes.Equal(data[0:2], expectedKeyboardChallengeResponseHeader) {
		return 0, fmt.Errorf("unexpected keyboard challenge response header: %v", data)
	}

	return binary.LittleEndian.Uint64(data[2:]), nil
}

// sendHostChallenge requests a challenge response from the device.
// Presumably this step exists to confirm it's a real Speed Editor.
func (ah AuthHandler) SendHostChallenge() error {
	_, err := ah.device.SendFeatureReport(featureReportHostChallenge)
	if err != nil {
		return fmt.Errorf("failed to send feature report: %w", err)
	}
	return nil
}

func (ah AuthHandler) GetHostChallengeResponse() ([]byte, error) {
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := ah.device.GetFeatureReport(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get feature report: %w", err)
	}

	if !bytes.Equal(data[0:2], expectedHostChallengeResponseHeader) {
		return nil, fmt.Errorf("unexpected host challenge response header: %v", data)
	}

	return data, nil
}

func (ah AuthHandler) SendAuthChallengeResponse(response uint64) error {
	responseBytes := make([]byte, len(authChallengeHeader))
	copy(responseBytes, authChallengeHeader)
	responseBytes = binary.LittleEndian.AppendUint64(authChallengeHeader, response)

	_, err := ah.device.SendFeatureReport(responseBytes)
	if err != nil {
		return fmt.Errorf("failed to send feature report: %w", err)
	}
	return nil
}

func (ah AuthHandler) GetAuthChallengeResult() (uint16, error) {
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := ah.device.GetFeatureReport(data)
	if err != nil {
		return 0, fmt.Errorf("failed to get feature report: %w", err)
	}

	if !bytes.Equal(data[0:2], expectedAuthResponseHeader) {
		return 0, fmt.Errorf("unexpected auth response header: %v", data)
	}

	return binary.LittleEndian.Uint16(data[2:4]), nil
}
