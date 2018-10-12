package main

import (
	"testing"
)

func TestLengthEncodedInt_1_byte(t *testing.T) {
	encoded := LengthEncodedInt(201)
	if encoded[0] != 0xC9 {
		t.Errorf("Bogus value for 1-byte encoded int: '%s'", encoded)
	}
}

func TestLengthEncodedInt_2_bytes(t *testing.T) {
	encoded := LengthEncodedInt(20103)
	if encoded[0] != 0x87 && encoded[1] != 0x4E {
		t.Errorf("Bogus value for 2-byte encoded int: '%s'", encoded)
	}
}

func TestLengthEncodedInt_3_bytes(t *testing.T) {
	encoded := LengthEncodedInt(2010378)
	if encoded[0] != 0x0A && encoded[1] != 0xAD && encoded[2] != 0x1E {
		t.Errorf("Bogus value for 3-byte encoded int: '%s'", encoded)
	}
}

func TestLengthEncodedInt_8_bytes(t *testing.T) {
	encoded := LengthEncodedInt(20103780195)
	// 04 ae 47 57 63
	if encoded[0] != 0x63 && encoded[1] != 0x57 && encoded[2] != 0x47 &&
		encoded[3] != 0xAE && encoded[4] != 0x04 && encoded[5] != 0x00 &&
		encoded[6] != 0x00 && encoded[7] != 0x00 {
		t.Errorf("Bogus value for 8-byte encoded int: '%s'", encoded)
	}
}
