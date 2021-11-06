package types

import (
	"testing"
)

func Test_DecodeIEEE11073float(t *testing.T) {
	cases := []struct {
		name  string
		data  []byte
		float float32
	}{
		{
			name:  "21.4",
			data:  []byte{0x00, 0x60, 0xA7, 0x20, 0xFB},
			float: 21.40,
		},
		{
			name:  "-21.4",
			data:  []byte{0x00, 0xA0, 0x58, 0xDF, 0xFB},
			float: -21.40,
		},
		{
			name:  "21.5",
			data:  []byte{0x00, 0xd7, 0x0, 0x0, 0xFF},
			float: 21.50,
		},
		{
			name:  "-21.5",
			data:  []byte{0x00, 0x29, 0xff, 0xff, 0xFf},
			float: -21.50,
		},
		{
			name:  "0.0",
			data:  []byte{0x00, 0x00, 0x00, 0x00, 0x00},
			float: 0.0,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			gotFloat := DecodeIEEE11073(tc.data)
			if gotFloat != tc.float {
				t.Errorf("decoding: got %.2f, expected %.2f", gotFloat, tc.float)
			}

			gotBytes := EncodeIEEE11073(tc.float)
			for i, b := range gotBytes {
				if b != tc.data[i] {
					t.Errorf("encoding: got %#v, expected %#v", gotBytes, tc.data)
					break
				}
			}
		})
	}
}
