package main

import (
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
		result = append(result, byte(num&0xFF), byte((num>>8)&0xFF))
	} else if num >= 65536 && num < 16777216 {
		result = append(result, byte(num&0xFF), byte((num>>8)&0xFF), byte((num>>16)&0xFF))
	} else {
		result = append(result,
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
