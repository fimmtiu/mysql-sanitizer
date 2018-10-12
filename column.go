package main

import (
	"fmt"
	"strings"
)

const TYPE_VARCHAR byte = 0x0F
const TYPE_TINY_BLOB byte = 0xF9
const TYPE_MEDIUM_BLOB byte = 0xFA
const TYPE_LONG_BLOB byte = 0xFB
const TYPE_BLOB byte = 0xFC
const TYPE_VAR_STRING byte = 0xFD
const TYPE_STRING byte = 0xFE

type Column struct {
	IsString bool
	Database string
	Table    string
	Name     string
	Length   uint32
}

func ReadColumn(parser *PacketParser) (Column, error) {
	column := Column{}
	parser.ReadVariableString() // catalog
	column.Database = strings.ToLower(parser.ReadVariableString())
	column.Table = strings.ToLower(parser.ReadVariableString())
	parser.ReadVariableString() // org_table
	column.Name = strings.ToLower(parser.ReadVariableString())
	parser.ReadVariableString() // org_name

	fixedFieldsLen := parser.ReadEncodedInt()
	if fixedFieldsLen != 12 {
		return Column{}, fmt.Errorf("Weird value for fixedFieldsLen: %d", fixedFieldsLen)
	}

	parser.ReadFixedInt2() // character set
	column.Length = parser.ReadFixedInt4()
	colType := parser.ReadFixedInt1()
	// We ignore the remaining fields

	if colType == TYPE_VARCHAR || colType == TYPE_TINY_BLOB || colType == TYPE_MEDIUM_BLOB ||
		colType == TYPE_LONG_BLOB || colType == TYPE_BLOB || colType == TYPE_VAR_STRING ||
		colType == TYPE_STRING {
		column.IsString = true
	} else {
		column.IsString = false
	}

	return column, nil
}

func (col Column) IsSafe() bool {
	// At this time, we believe that all non-string columns are safe.
	if !col.IsString {
		return true
	}

	// Don't mangle EXPLAIN and DESCRIBE statements.
	if col.Database == "information_schema" && (col.Table == "columns" || col.Table == "schemata") {
		return true
	}

	// Allow viewing the values of internal stuff like EXPLAIN output and "@@" MySQL variables.
	if col.Database == "" && col.Table == "" {
		return true
	}

	// If we've explicitly permitted this column in the JSON list, it's safe.
	if whitelist.IsColumnPresent(col.Database, col.Table, col.Name) {
		return true
	}

	return false
}
