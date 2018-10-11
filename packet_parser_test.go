package main

import (
	"testing"

	"github.com/pubnative/mysqlproto-go"
)

func TestReadEncodedInt_1(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("\x21")}
	parser := NewPacketParser(packet)
	value := parser.ReadEncodedInt()
	if value != 33 {
		t.Errorf("Bogus value for 1-byte encoded int: 0x%02x", value)
	}
	if parser.offset != 1 {
		t.Errorf("Bogus value for offset after 1-byte int read: 0x%02x", parser.offset)
	}
}

func TestReadEncodedInt_2(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("\xFC\x21\x01")}
	parser := NewPacketParser(packet)
	value := parser.ReadEncodedInt()
	if value != 289 {
		t.Errorf("Bogus value for 2-byte encoded int: 0x%02x", value)
	}
	if parser.offset != 3 {
		t.Errorf("Bogus value for offset after 2-byte int read: 0x%02x", parser.offset)
	}
}

func TestReadEncodedInt_3(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("\xFD\x21\x01\x03")}
	parser := NewPacketParser(packet)
	value := parser.ReadEncodedInt()
	if value != 196897 {
		t.Errorf("Bogus value for 8-byte encoded int: 0x%02x", value)
	}
	if parser.offset != 4 {
		t.Errorf("Bogus value for offset after 3-byte int read: 0x%02x", parser.offset)
	}
}

func TestReadEncodedInt_8(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("\xFE\x21\x01\x03\x68\x00\x00\x00\x00\x00")}
	parser := NewPacketParser(packet)
	value := parser.ReadEncodedInt()
	if value != 1745027361 {
		t.Errorf("Bogus value for 8-byte encoded int: 0x%02x", value)
	}
	if parser.offset != 9 {
		t.Errorf("Bogus value for offset after 8-byte int read: 0x%02x", parser.offset)
	}
}

func TestReadFixedString(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("hello world!")}
	parser := NewPacketParser(packet)
	value := parser.ReadFixedString(7)
	if value != "hello w" {
		t.Errorf("Bogus value for variable string: '%s'", value)
	}
	if parser.offset != 7 {
		t.Errorf("Bogus value for offset after fixed string read: 0x%02x", parser.offset)
	}
}

func TestReadNullTermString(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("hello world!\x00foo")}
	parser := NewPacketParser(packet)
	value := parser.ReadNullTermString()
	if value != "hello world!" {
		t.Errorf("Bogus value for variable string: '%s'", value)
	}
	if parser.offset != 13 {
		t.Errorf("Bogus value for offset after null-terminated string read: 0x%02x", parser.offset)
	}
}

func TestReadVariableString(t *testing.T) {
	packet := &mysqlproto.Packet{0, []byte("\x0chello world!")}
	parser := NewPacketParser(packet)
	value := parser.ReadVariableString()
	if value != "hello world!" {
		t.Errorf("Bogus value for variable string: '%s'", value)
	}
	if parser.offset != 13 {
		t.Errorf("Bogus value for offset after variable string read: 0x%02x", parser.offset)
	}
}
