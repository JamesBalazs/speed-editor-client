// Based on the authentication algorithm by Sylvain Munaut <tnt@246tNt.com> (Apache-2.0)
// https://github.com/smunaut/blackmagic-misc/blob/master/bmd.py
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/sstallion/go-hid"
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

func Authenticate(d *hid.Device) {
	resetAuthState(d)
	challenge := getKeyboardChallenge(d)

	sendHostChallenge(d)
	_ = getHostChallengeResponse(d) // We don't care about the response, since we don't care if it's a real Speed Editor

	response := calculateChallengeResponse(challenge)
	sendAuthChallengeResponse(d, response)

	reauthTimeout = getAuthChallengeResult(d)

	fmt.Printf("Auth success! Reauth in %d\n", reauthTimeout)
}

func resetAuthState(d *hid.Device) {
	_, err := d.SendFeatureReport(featureReportDefaultState)
	if err != nil {
		panic(err.Error())
	}
}

func getKeyboardChallenge(d *hid.Device) uint64 {
	// get keyboard challenge, store in a new copy of the byte array
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := d.GetFeatureReport(data)
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
func sendHostChallenge(d *hid.Device) {
	_, err := d.SendFeatureReport(featureReportHostChallenge)
	if err != nil {
		panic(err.Error())
	}
}

func getHostChallengeResponse(d *hid.Device) []byte {
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := d.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedHostChallengeResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}

	return data
}

func sendAuthChallengeResponse(d *hid.Device, response uint64) {
	responseBytes := make([]byte, len(authChallengeHeader))
	copy(responseBytes, authChallengeHeader)
	responseBytes = binary.LittleEndian.AppendUint64(authChallengeHeader, response)

	_, err := d.SendFeatureReport(responseBytes)
	if err != nil {
		panic(err.Error())
	}
}

func getAuthChallengeResult(d *hid.Device) uint16 {
	data := make([]byte, len(featureReportDefaultState))
	copy(data, featureReportDefaultState)

	_, err := d.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedAuthResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}

	return binary.LittleEndian.Uint16(data[2:4])
}
