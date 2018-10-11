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
	output.Dump(contents, "Writing packet:\n")
	stream.Write(contents)
}
