package main

import (
	"fmt"

	"github.com/pubnative/mysqlproto-go"
)

type PacketParser struct {
	data   []byte
	offset uint64
}

func NewPacketParser(packet mysqlproto.Packet) *PacketParser {
	return &PacketParser{packet.Payload, 0}
}

func (parser *PacketParser) ReadEncodedInt() uint64 {
	if parser.data[parser.offset] < 0xFB {
		return uint64(parser.ReadFixedInt1())
	} else if parser.data[parser.offset] == 0xFC {
		parser.offset++
		return uint64(parser.ReadFixedInt2())
	} else if parser.data[parser.offset] == 0xFD {
		parser.offset++
		return uint64(parser.ReadFixedInt3())
	} else if parser.data[parser.offset] == 0xFE {
		parser.offset++
		return uint64(parser.ReadFixedInt8())
	} else {
		panic(fmt.Sprintf("Invalid header byte for length-encoded integer: 0x%02x!", parser.data[0]))
	}
}

func (parser *PacketParser) ReadFixedString(length uint64) string {
	bytes := parser.data[parser.offset : parser.offset+length]
	parser.offset += length
	return string(bytes)
}

func (parser *PacketParser) ReadNullTermString() string {
	var null_index int = -1

	for i, byte := range parser.data[parser.offset:] {
		if byte == 0x00 {
			null_index = i
			break
		}
	}

	if null_index < 0 {
		panic("Didn't find a NUL when looking for a null-terminated string!")
	}

	bytes := parser.data[parser.offset : parser.offset+uint64(null_index)]
	parser.offset += uint64(len(bytes)) + 1
	return string(bytes)
}

func (parser *PacketParser) ReadVariableString() string {
	strlen := parser.ReadEncodedInt()
	bytes := parser.data[parser.offset : parser.offset+strlen]
	parser.offset += strlen
	return string(bytes)
}

func (parser *PacketParser) ReadFixedInt1() uint8 {
	fixedInt := parser.data[parser.offset]
	parser.offset++
	return fixedInt
}

func (parser *PacketParser) ReadFixedInt2() uint16 {
	var fixedInt uint16 = uint16(parser.data[parser.offset+1])<<8 |
		uint16(parser.data[parser.offset])
	parser.offset += 2
	return fixedInt
}

func (parser *PacketParser) ReadFixedInt3() uint32 {
	var fixedInt uint32 = uint32(parser.data[parser.offset+2])<<16 |
		uint32(parser.data[parser.offset+1])<<8 |
		uint32(parser.data[parser.offset])
	parser.offset += 3
	return fixedInt
}

func (parser *PacketParser) ReadFixedInt4() uint32 {
	var fixedInt uint32 = uint32(parser.data[parser.offset+3])<<24 |
		uint32(parser.data[parser.offset+2])<<16 |
		uint32(parser.data[parser.offset+1])<<8 |
		uint32(parser.data[parser.offset])
	parser.offset += 4
	return fixedInt
}

// And of course I can't line this up nicely, because gofmt.
func (parser *PacketParser) ReadFixedInt8() uint64 {
	var fixedInt uint64 = uint64(parser.data[parser.offset+7])<<56 |
		uint64(parser.data[parser.offset+6])<<48 |
		uint64(parser.data[parser.offset+5])<<40 |
		uint64(parser.data[parser.offset+4])<<32 |
		uint64(parser.data[parser.offset+3])<<24 |
		uint64(parser.data[parser.offset+2])<<16 |
		uint64(parser.data[parser.offset+1])<<8 |
		uint64(parser.data[parser.offset])
	parser.offset += 8
	return fixedInt
}
