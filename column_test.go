package main

import (
	"testing"

	"github.com/pubnative/mysqlproto-go"
)

func TestReadColumn(t *testing.T) {
	packet := mysqlproto.Packet{0, []byte("\x03def\x04honk\x04bonk\x04bonk\x08woopwoop\x08woopwoop\x0c\x21\x00\xfd\x02\x00\x00\xfd\x00\x00\x00\x00\x00")}
	parser := NewPacketParser(packet)
	column, err := ReadColumn(parser)
	if err != nil {
		t.Errorf("ReadColumn failed: %s", err)
	}
	if !column.IsString {
		t.Error("Column should be a string!")
	}
	if column.Database != "honk" {
		t.Errorf("Unexpected value for column.Database: '%s'", column.Database)
	}
	if column.Table != "bonk" {
		t.Errorf("Unexpected value for column.Table: '%s'", column.Table)
	}
	if column.Name != "woopwoop" {
		t.Errorf("Unexpected value for column.Name: '%s'", column.Name)
	}
}
