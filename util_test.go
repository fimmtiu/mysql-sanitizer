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
	if encoded[0] != 0xFC || encoded[1] != 0x87 || encoded[2] != 0x4E {
		t.Errorf("Bogus value for 2-byte encoded int: '%s'", encoded)
	}
}

func TestLengthEncodedInt_3_bytes(t *testing.T) {
	encoded := LengthEncodedInt(2010378)
	if encoded[0] != 0xFD || encoded[1] != 0x0A || encoded[2] != 0xAD || encoded[3] != 0x1E {
		t.Errorf("Bogus value for 3-byte encoded int: '%s'", encoded)
	}
}

func TestLengthEncodedInt_8_bytes(t *testing.T) {
	encoded := LengthEncodedInt(20103780195)
	if encoded[0] != 0xFE || encoded[1] != 0x63 || encoded[2] != 0x57 ||
		encoded[3] != 0x47 || encoded[4] != 0xAE || encoded[5] != 0x04 ||
		encoded[6] != 0x00 || encoded[7] != 0x00 || encoded[8] != 0x00 {
		t.Errorf("Bogus value for 8-byte encoded int: '%s'", encoded)
	}
}

func TestVariableString_1_byte(t *testing.T) {
	bytes := VariableString("I've got %d legs", 2)
	if string(bytes) != "\x0fI've got 2 legs" {
		t.Errorf("Unexpected result from VariableString: '%s'", string(bytes))
	}
}

func TestVariableString_2_bytes(t *testing.T) {
	bytes := VariableString("Lo, praise of the prowess of people-kings of spear-armed Danes, in days long sped, we have heard, and what honor the athelings won! Oft Scyld the Scefing from squadroned foes, from many a tribe, the mead-bench tore, awing the earls. Since erst he lay friendless, a foundling, fate repaid him: for he waxed under welkin, in wealth he throve, till before him the folk, both far and near, who house by the whale-path, heard his mandate, gave him gifts: a good king he!")
	if string(bytes) != "\xfc\xd1\x01Lo, praise of the prowess of people-kings of spear-armed Danes, in days long sped, we have heard, and what honor the athelings won! Oft Scyld the Scefing from squadroned foes, from many a tribe, the mead-bench tore, awing the earls. Since erst he lay friendless, a foundling, fate repaid him: for he waxed under welkin, in wealth he throve, till before him the folk, both far and near, who house by the whale-path, heard his mandate, gave him gifts: a good king he!" {
		output.Log("Initial bytes: 0x%02x 0x%02x 0x%02x 0x%02x", bytes[0], bytes[1], bytes[2], bytes[3])
		t.Errorf("Unexpected result from VariableString: '%s'", string(bytes))
	}
}

func TestErrorPacket(t *testing.T) {
	packet := ErrorPacket(32, 31337, "HONK1", "Your packet makes me sad.")
	if packet.Payload[0] != 0xFF {
		t.Errorf("Packet doesn't start with a ERR header: 0x%02x.", packet.Payload[0])
	}
	if packet.SequenceID != 33 {
		t.Errorf("Sequence ID isn't one more than the original: %d.", packet.SequenceID)
	}
	if packet.Payload[1] != 0x69 || packet.Payload[2] != 0x7A {
		t.Errorf("Unexpected error code: 0x%02x%02x.", packet.Payload[2], packet.Payload[1])
	}
	if packet.Payload[3] != 0x23 {
		t.Errorf("Didn't see the SQL state marker: 0x%02x.", packet.Payload[3])
	}
	if string(packet.Payload[4:9]) != "HONK1" {
		t.Errorf("Incorrect SQL state marker: '%s'", string(packet.Payload[4:9]))
	}
}
