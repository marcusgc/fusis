// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: jsing@google.com (Joel Sing)

package healthcheck

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

var (
	testAuthenticator1 = &radiusAuthenticator{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}
	testAuthenticator2 = &radiusAuthenticator{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
)

var radiusPasswordTests = []struct {
	passwd        string
	secret        string
	authenticator *radiusAuthenticator
	want          []byte
}{
	{
		"",
		"",
		testAuthenticator1,
		[]byte{},
	},
	{
		"a",
		"secret",
		testAuthenticator1,
		[]byte{
			0x37, 0x0e, 0xc9, 0x44, 0x90, 0x92, 0xd3, 0x67,
			0xca, 0x5e, 0xfb, 0x70, 0x6c, 0xd4, 0xe7, 0x07,
		},
	},
	{
		"0123456789abcdef",
		"mySuperRADIUSSecret",
		testAuthenticator1,
		[]byte{
			0xca, 0x9d, 0xb6, 0x64, 0xec, 0x31, 0xd8, 0xac,
			0x40, 0xe8, 0x85, 0x90, 0x40, 0x5d, 0x92, 0xd6,
		},
	},
	{
		"0123456789abcdef",
		"mySuperRADIUSSecret!",
		testAuthenticator1,
		[]byte{
			0x7d, 0x8f, 0x1f, 0x9f, 0x31, 0xb9, 0xb4, 0xe9,
			0x30, 0x56, 0xe7, 0x63, 0x7c, 0xa4, 0x82, 0x4d,
		},
	},
	{
		"0123456789abcdef!",
		"mySuperRADIUSSecret!",
		testAuthenticator1,
		[]byte{
			0x7d, 0x8f, 0x1f, 0x9f, 0x31, 0xb9, 0xb4, 0xe9,
			0x30, 0x56, 0xe7, 0x63, 0x7c, 0xa4, 0x82, 0x4d,
			0x5f, 0x0b, 0xd0, 0x56, 0x5b, 0x1a, 0x1f, 0x64,
			0xc5, 0x78, 0xb2, 0x69, 0x48, 0x8a, 0x7d, 0x47,
		},
	},
	{
		"0123456789abcdef!",
		"mySuperRADIUSSecret!",
		testAuthenticator2,
		[]byte{
			0xd8, 0xbd, 0x08, 0x45, 0x9a, 0x3a, 0x9d, 0x5c,
			0xc4, 0x3e, 0xf0, 0x49, 0x73, 0x2f, 0xe4, 0xfd,
			0xa6, 0xd1, 0xb8, 0x68, 0x07, 0x7c, 0xe3, 0x74,
			0x9b, 0x82, 0xd2, 0x31, 0x7c, 0xf5, 0x71, 0x38,
		},
	},
	{
		strings.Repeat("a", 256),
		"mySuperRADIUSSecret!",
		testAuthenticator1,
		[]byte{
			0x2c, 0xdf, 0x4c, 0xcd, 0x64, 0xed, 0xe3, 0xbf,
			0x69, 0x0e, 0xe7, 0x60, 0x7e, 0xa1, 0x86, 0x4a,
			0xbf, 0x82, 0xc2, 0xd8, 0x15, 0x02, 0x63, 0x11,
			0x1c, 0x00, 0x73, 0x0d, 0xf1, 0xb3, 0xb6, 0x93,
			0x83, 0x54, 0xcf, 0x94, 0x2b, 0xfc, 0x8a, 0x9d,
			0xca, 0x57, 0xd8, 0x5b, 0x7c, 0x8b, 0xdb, 0x2a,
			0x01, 0x65, 0x0d, 0x74, 0x80, 0x54, 0xd7, 0x4a,
			0x0a, 0x95, 0x60, 0xc5, 0x4d, 0x76, 0x2a, 0x5e,
			0x37, 0x0b, 0x87, 0xc2, 0x60, 0x89, 0x90, 0x46,
			0x04, 0x66, 0x61, 0x8e, 0x7e, 0xc6, 0x55, 0x4f,
			0x94, 0x73, 0xb9, 0xe2, 0xfe, 0x9f, 0xf6, 0x62,
			0x60, 0x9d, 0xf5, 0xf3, 0x37, 0xc0, 0xfa, 0x74,
			0x5c, 0x0a, 0x21, 0x5d, 0x94, 0x4d, 0x30, 0x2a,
			0x03, 0xe3, 0xb9, 0x69, 0xd9, 0x1f, 0xe3, 0xf8,
			0x09, 0x20, 0x4e, 0xc3, 0x0a, 0xb9, 0x89, 0xef,
			0x7c, 0x49, 0x8c, 0xf5, 0x24, 0x38, 0x93, 0x2a,
		},
	},
}

func TestRADIUSPassword(t *testing.T) {
	for i, rt := range radiusPasswordTests {
		got := radiusPassword(rt.passwd, rt.secret, rt.authenticator)
		if !bytes.Equal(got, rt.want) {
			t.Errorf("Test %d: got password %#v, want %#v", i, got, rt.want)
		}
	}
}

var radiusResponseAuthenticatorTests = []struct {
	packet               []byte
	secret               string
	requestAuthenticator *radiusAuthenticator
}{
	{
		[]byte{
			0x03, 0x76, 0x00, 0x33, 0x89, 0xf1, 0xd6, 0xe2,
			0xbb, 0xfd, 0x53, 0x37, 0x16, 0x64, 0x90, 0xe2,
			0x38, 0x20, 0xcd, 0x26, 0x12, 0x1f, 0x72, 0x61,
			0x64, 0x69, 0x75, 0x73, 0x31, 0x2d, 0x34, 0x2e,
			0x74, 0x77, 0x64, 0x2e, 0x63, 0x6f, 0x72, 0x70,
			0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
			0x63, 0x6f, 0x6d,
		},
		"mySuperRADIUSSecret!",
		testAuthenticator1,
	},
	{
		[]byte{
			0x03, 0x76, 0x00, 0x33, 0xd4, 0x2c, 0xbe, 0xcd,
			0xb1, 0x87, 0x1f, 0x0d, 0x47, 0x61, 0x62, 0xae,
			0x35, 0xbd, 0xd3, 0x9e, 0x12, 0x1f, 0x72, 0x61,
			0x64, 0x69, 0x75, 0x73, 0x31, 0x2d, 0x34, 0x2e,
			0x74, 0x77, 0x64, 0x2e, 0x63, 0x6f, 0x72, 0x70,
			0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
			0x63, 0x6f, 0x6d,
		},
		"mySuperRADIUSSecret!",
		testAuthenticator2,
	},
}

func TestRADIUSResponseAuthenticator(t *testing.T) {
	for i, rt := range radiusResponseAuthenticatorTests {
		b := bytes.NewReader(rt.packet)
		rp := &radiusPacket{}
		if err := rp.decode(b); err != nil {
			t.Errorf("Test %d: failed to decode RADIUS packet: %v", i, err)
			continue
		}
		got, err := responseAuthenticator(rp, rt.requestAuthenticator, rt.secret)
		if err != nil {
			t.Errorf("Test %d: failed to calculate response authenticator: %v", i, err)
			continue
		}
		want := &rp.Authenticator
		if !bytes.Equal(got[:], want[:]) {
			t.Errorf("Test %d: got response authenticator %#v, want %#v", i, got, want)
		}
	}
}

var radiusPacketTests = []struct {
	packet       []byte
	radiusPacket *radiusPacket
}{
	{
		[]byte{
			0x01, 0x13, 0x00, 0x65, 0x07, 0x5a, 0x3d, 0xaa,
			0xe3, 0x67, 0xb4, 0xb0, 0x3b, 0xdb, 0x03, 0xcf,
			0xb1, 0x11, 0x0c, 0xa7, 0x20, 0x0e, 0x72, 0x61,
			0x64, 0x69, 0x75, 0x73, 0x70, 0x72, 0x6f, 0x62,
			0x65, 0x72, 0x01, 0x0f, 0x72, 0x61, 0x64, 0x69,
			0x75, 0x73, 0x2d, 0x70, 0x72, 0x6f, 0x62, 0x65,
			0x72, 0x02, 0x22, 0x85, 0x71, 0xb2, 0x4f, 0x7e,
			0x81, 0x54, 0x59, 0x5a, 0x12, 0x58, 0xb0, 0x71,
			0x93, 0xb7, 0x36, 0xc9, 0xdd, 0x65, 0x2f, 0x6f,
			0x07, 0x87, 0xe9, 0xce, 0x34, 0x8b, 0xc4, 0xef,
			0xf6, 0x19, 0xfd, 0x04, 0x06, 0xac, 0x17, 0x00,
			0x08, 0x05, 0x06, 0x00, 0x00, 0x00, 0x2a, 0x06,
			0x06, 0x00, 0x00, 0x00, 0x01,
		},
		&radiusPacket{
			radiusHeader{
				Code:       0x1,
				Identifier: 0x13,
				Length:     0x65,
				Authenticator: radiusAuthenticator{
					0x07, 0x5a, 0x3d, 0xaa, 0xe3, 0x67, 0xb4, 0xb0,
					0x3b, 0xdb, 0x03, 0xcf, 0xb1, 0x11, 0x0c, 0xa7,
				},
			},
			[]*radiusAttribute{
				{
					raType: 0x20,
					length: 0x0e,
					value: []byte{
						0x72, 0x61, 0x64, 0x69, 0x75, 0x73, 0x70, 0x72,
						0x6f, 0x62, 0x65, 0x72,
					},
				},
				{
					raType: 0x1,
					length: 0xf,
					value: []byte{
						0x72, 0x61, 0x64, 0x69, 0x75, 0x73, 0x2d, 0x70,
						0x72, 0x6f, 0x62, 0x65, 0x72,
					},
				},
				{
					raType: 0x2,
					length: 0x22,
					value: []byte{
						0x85, 0x71, 0xb2, 0x4f, 0x7e, 0x81, 0x54, 0x59,
						0x5a, 0x12, 0x58, 0xb0, 0x71, 0x93, 0xb7, 0x36,
						0xc9, 0xdd, 0x65, 0x2f, 0x6f, 0x07, 0x87, 0xe9,
						0xce, 0x34, 0x8b, 0xc4, 0xef, 0xf6, 0x19, 0xfd,
					},
				},
				{
					raType: 0x4,
					length: 0x6,
					value:  []byte{0xac, 0x17, 0x00, 0x08},
				},
				{
					raType: 0x5,
					length: 0x6,
					value:  []byte{0x00, 0x00, 0x00, 0x2a},
				},
				{
					raType: 0x6,
					length: 0x6,
					value:  []byte{0x00, 0x00, 0x00, 0x01},
				},
			},
		},
	},
	{
		[]byte{
			0x03, 0x76, 0x00, 0x33, 0x12, 0x35, 0xf4, 0xf2,
			0xb8, 0xdc, 0x32, 0xda, 0x0d, 0x6a, 0x3f, 0x3d,
			0xed, 0xbd, 0x78, 0x95, 0x12, 0x1f, 0x72, 0x61,
			0x64, 0x69, 0x75, 0x73, 0x31, 0x2d, 0x34, 0x2e,
			0x74, 0x77, 0x64, 0x2e, 0x63, 0x6f, 0x72, 0x70,
			0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
			0x63, 0x6f, 0x6d,
		},
		&radiusPacket{
			radiusHeader{
				Code:       0x3,
				Identifier: 0x76,
				Length:     0x33,
				Authenticator: radiusAuthenticator{
					0x12, 0x35, 0xf4, 0xf2, 0xb8, 0xdc, 0x32, 0xda,
					0x0d, 0x6a, 0x3f, 0x3d, 0xed, 0xbd, 0x78, 0x95,
				},
			},
			[]*radiusAttribute{
				{
					raType: 0x12,
					length: 0x1f,
					value: []byte{
						0x72, 0x61, 0x64, 0x69, 0x75, 0x73, 0x31, 0x2d,
						0x34, 0x2e, 0x74, 0x77, 0x64, 0x2e, 0x63, 0x6f,
						0x72, 0x70, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
						0x65, 0x2e, 0x63, 0x6f, 0x6d,
					},
				},
			},
		},
	},
}

func TestRADIUSPacket(t *testing.T) {
	for i, rt := range radiusPacketTests {
		b := bytes.NewReader(rt.packet)
		rp := &radiusPacket{}
		if err := rp.decode(b); err != nil {
			t.Errorf("Test %d: failed to decode packet: %v", i, err)
			continue
		}
		if !reflect.DeepEqual(rp, rt.radiusPacket) {
			t.Errorf("Test %d: got RADIUS packet %#v, want %#v", i, rp, rt.radiusPacket)
		}
	}
}