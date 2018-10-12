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
	Alias    string
	Name     string
	Length   uint32
}

func ReadColumn(parser *PacketParser) (Column, error) {
	column := Column{}
	parser.ReadVariableString() // catalog
	column.Database = strings.ToLower(parser.ReadVariableString())
	parser.ReadVariableString()                                 // possibly aliased table name, ignore
	column.Table = strings.ToLower(parser.ReadVariableString()) // real table name
	column.Alias = parser.ReadVariableString()                  // possibly aliased column name
	column.Name = strings.ToLower(parser.ReadVariableString())  // real column name

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

	output.Log("Column: database '%s', table '%s', name '%s' ('%s')", column.Database, column.Table, column.Name, column.Alias)
	return column, nil
}

func (col Column) IsSafe() bool {
	// At this time, we believe that all non-string columns are safe.
	if !col.IsString {
		return true
	}

	// Don't mangle EXPLAIN and DESCRIBE statements.
	if col.Database == "information_schema" &&
		(col.Table == "columns" || col.Table == "schemata" || col.Table == "table_names") {
		return true
	}

	// Allow viewing the values of internal stuff like EXPLAIN output and "@@" MySQL variables.
	if col.Database == "" && col.Table == "" {
		// But not things like `CONCAT(address, " ")`, which suuuuuck. (We
		// should inspect those more closely, but that's for later when we
		// actually start parsing SQL.)
		if col.Name == "" {
			return !strings.Contains(col.Name, "(")
		} else {
			return !strings.Contains(col.Alias, "(")
		}
	}

	// If we've explicitly permitted this column in the JSON list, it's safe.
	if whitelist.IsColumnPresent(col.Database, col.Table, col.Name) {
		return true
	}

	return false
}
