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
		return uint64(parser.readFixedInt1(parser.offset))
	} else if parser.data[parser.offset] == 0xFC {
		return uint64(parser.readFixedInt2(parser.offset + 1))
	} else if parser.data[parser.offset] == 0xFD {
		return uint64(parser.readFixedInt3(parser.offset + 1))
	} else if parser.data[parser.offset] == 0xFE {
		return uint64(parser.readFixedInt8(parser.offset + 1))
	} else {
		panic(fmt.Sprintf("Invalid header byte for length-encoded integer: 0x%02x!", parser.data[0]))
	}
}

func (parser *PacketParser) ReadVariableString() string {
	strlen := parser.ReadEncodedInt()
	bytes := parser.data[parser.offset : strlen+1]
	parser.offset += strlen
	return string(bytes)
}

func (parser *PacketParser) readFixedInt1(offset uint64) uint8 {
	fixedInt := parser.data[offset]
	parser.offset = offset + 1
	return fixedInt
}

func (parser *PacketParser) readFixedInt2(offset uint64) uint16 {
	var fixedInt uint16 = uint16(parser.data[offset+1])<<8 | uint16(parser.data[offset])
	parser.offset = offset + 2
	return fixedInt
}

func (parser *PacketParser) readFixedInt3(offset uint64) uint32 {
	var fixedInt uint32 = uint32(parser.data[offset+2])<<16 | uint32(parser.data[offset+1])<<8 | uint32(parser.data[offset])
	parser.offset = offset + 3
	return fixedInt
}

// And of course I can't line this up nicely, because gofmt.
func (parser *PacketParser) readFixedInt8(offset uint64) uint64 {
	var fixedInt uint64 = uint64(parser.data[offset+7])<<56 |
		uint64(parser.data[offset+6])<<48 |
		uint64(parser.data[offset+5])<<40 |
		uint64(parser.data[offset+4])<<32 |
		uint64(parser.data[offset+3])<<24 |
		uint64(parser.data[offset+2])<<16 |
		uint64(parser.data[offset+1])<<8 |
		uint64(parser.data[offset])
	parser.offset = offset + 8
	return fixedInt
}
