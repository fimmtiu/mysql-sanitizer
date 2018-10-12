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

func TestColumnIsSafe_NotString(t *testing.T) {
	column := Column{false, "honk", "bonk", "woopwoop", 255}
	if !column.IsSafe() {
		t.Error("Non-string columns should always be safe!")
	}
}

func TestColumnIsSafe_String(t *testing.T) {
	column := Column{true, "honk", "bonk", "woopwoop", 255}
	if column.IsSafe() {
		t.Error("Non-whitelisted string columns shouldn't be safe!")
	}
}

func TestColumnIsSafe_InfoSchema(t *testing.T) {
	column := Column{true, "information_schema", "columns", "woopwoop", 255}
	if !column.IsSafe() {
		t.Error("information_schema.columns should always be safe!")
	}

	column = Column{true, "information_schema", "schemata", "woopwoop", 255}
	if !column.IsSafe() {
		t.Error("information_schema.schemata should always be safe!")
	}

	column = Column{true, "information_schema", "table_names", "woopwoop", 255}
	if !column.IsSafe() {
		t.Error("information_schema.table_names should always be safe!")
	}

	column = Column{true, "information_schema", "user_privileges", "woopwoop", 255}
	if column.IsSafe() {
		t.Error("Other information_schema tables aren't safe!")
	}
}

func TestColumnIsSafe_Internals(t *testing.T) {
	column := Column{true, "", "", "@@woopwoop", 255}
	if !column.IsSafe() {
		t.Error("Columns without a schema should always be safe!")
	}
}
