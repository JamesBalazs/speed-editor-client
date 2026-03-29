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
	Authenticate() time.Duration
	ResetAuthState()
	GetKeyboardChallenge() uint64
	SendHostChallenge()
	GetHostChallengeResponse() []byte
	SendAuthChallengeResponse(response uint64)
	GetAuthChallengeResult() uint16
}

type AuthHandler struct {
	device deviceInterface
}

// Authenticate handles the entire handshake between the host and the Speed Editor.
//
// It returns the duration before the Speed Editor expects a reauth.
func (ah AuthHandler) Authenticate() time.Duration {
	ah.ResetAuthState()
	challenge := ah.GetKeyboardChallenge()

	ah.SendHostChallenge()
	_ = ah.GetHostChallengeResponse() // We don't care about the response, since we don't care if it's a real Speed Editor

	response := auth.CalculateChallengeResponse(challenge)
	ah.SendAuthChallengeResponse(response)

	reauthTimeout = ah.GetAuthChallengeResult()
	return time.Duration(reauthTimeout-10) * time.Second
}

func (ah AuthHandler) ResetAuthState() {
	_, err := ah.device.SendFeatureReport(featureReportDefaultState)
	if err != nil {
		panic(err.Error())
	}
}

func (ah AuthHandler) GetKeyboardChallenge() uint64 {
	// get keyboard challenge, store in a new copy of the byte array
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := ah.device.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedKeyboardChallengeResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}

	return binary.LittleEndian.Uint64(data[2:])
}

// sendHostChallenge requests a challenge response from the device.
// Presumably this step exists to confirm it's a real Speed Editor.
func (ah AuthHandler) SendHostChallenge() {
	_, err := ah.device.SendFeatureReport(featureReportHostChallenge)
	if err != nil {
		panic(err.Error())
	}
}

func (ah AuthHandler) GetHostChallengeResponse() []byte {
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := ah.device.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedHostChallengeResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}

	return data
}

func (ah AuthHandler) SendAuthChallengeResponse(response uint64) {
	responseBytes := make([]byte, len(authChallengeHeader))
	copy(responseBytes, authChallengeHeader)
	responseBytes = binary.LittleEndian.AppendUint64(authChallengeHeader, response)

	_, err := ah.device.SendFeatureReport(responseBytes)
	if err != nil {
		panic(err.Error())
	}
}

func (ah AuthHandler) GetAuthChallengeResult() uint16 {
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := ah.device.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedAuthResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}

	return binary.LittleEndian.Uint16(data[2:4])
}
