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

	expectedKeyboardChallengeResponseHeader = []byte{0x06, 0x00}
	expectedHostChallengeResponseHeader     = []byte{0x06, 0x02}
)

func Authenticate(d *hid.Device) {
	resetAuthState(d)
	challenge := getKeyboardChallenge(d)
	sendHostChallenge(d)

	_ = challenge

	fmt.Print("Auth success!")
}

func resetAuthState(d *hid.Device) {
	_, err := d.SendFeatureReport(featureReportDefaultState)
	if err != nil {
		panic(err.Error())
	}
}

func getKeyboardChallenge(d *hid.Device) uint64 {
	data := featureReportDefaultState // get keyboard challenge, store in a new copy of the byte array
	_, err := d.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedKeyboardChallengeResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}

	return binary.LittleEndian.Uint64(data[2:])
}

// sendHostChallenge requests a challenge response from the device, presumably this step exists
// to confirm it's a real Speed Editor.
//
// We don't care about the response, since we don't care if it's a real Speed Editor
func sendHostChallenge(d *hid.Device) {
	_, err := d.SendFeatureReport(featureReportHostChallenge)
	if err != nil {
		panic(err.Error())
	}

	data := featureReportDefaultState
	_, err = d.GetFeatureReport(data)
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(data[0:2], expectedHostChallengeResponseHeader) {
		panic(fmt.Sprintf("Unexpected auth response header: %v", data))
	}
}
