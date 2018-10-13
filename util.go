package main

import (
	"fmt"

	"github.com/pubnative/mysqlproto-go"
)

func WritePacket(stream *mysqlproto.Stream, packet mysqlproto.Packet) {
	contents := make([]byte, len(packet.Payload)+4)
	contents[0] = byte(len(packet.Payload) & 0xFF)
	contents[1] = byte((len(packet.Payload) >> 8) & 0xFF)
	contents[2] = byte((len(packet.Payload) >> 16) & 0xFF)
	contents[3] = packet.SequenceID
	copied := copy(contents[4:], packet.Payload)
	if copied != len(packet.Payload) {
		panic("wtf")
	}
	stream.Write(contents)
}

func LengthEncodedInt(num uint) []byte {
	result := make([]byte, 0, 8)

	if num < 251 {
		result = append(result, byte(num))
	} else if num >= 251 && num < 65536 {
		result = append(result, 0xFC, byte(num&0xFF), byte((num>>8)&0xFF))
	} else if num >= 65536 && num < 16777216 {
		result = append(result, 0xFD, byte(num&0xFF), byte((num>>8)&0xFF), byte((num>>16)&0xFF))
	} else {
		result = append(result, 0xFE,
			byte(num&0xFF),
			byte((num>>8)&0xFF),
			byte((num>>16)&0xFF),
			byte((num>>24)&0xFF),
			byte((num>>32)&0xFF),
			byte((num>>40)&0xFF),
			byte((num>>48)&0xFF),
			byte((num>>56)&0xFF),
		)
	}

	return result
}

func VariableString(format string, args ...interface{}) []byte {
	str := fmt.Sprintf(format, args...)
	length := LengthEncodedInt(uint(len(str)))
	return append(length, []byte(str)...)
}

func ErrorPacket(sequenceId byte, code int, sqlState string, format string, args ...interface{}) mysqlproto.Packet {
	// Screw it! We'll assume that everyone has CLIENT_PROTOCOL_41. It's
	// only been 14 years since it came out...
	chunks := []byte{0xFF}
	chunks = append(chunks, byte(code&0xFF), byte((code>>8)&0xFF))
	chunks = append(chunks, 0x23) // There's no documentation anywhere about what this field is.
	if len(sqlState) > 5 {
		panic(fmt.Sprintf("Bogus SQL state for ErrorPacket: '%s'", sqlState))
	}
	chunks = append(chunks, []byte(sqlState)...)
	str := fmt.Sprintf(format, args...)
	chunks = append(chunks, []byte(str)...)

	return mysqlproto.Packet{sequenceId + 1, chunks}
}
