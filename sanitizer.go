package main

import (
	"github.com/pubnative/mysqlproto-go"
)

type Sanitizer struct {
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{}
}

func (sanitizer *Sanitizer) Sanitize(packet mysqlproto.Packet) mysqlproto.Packet {
	// if packet.Payload[0] != 0 && FIXME whatever  {
	// 	// We only need to sanitize packets that contain response data.
	// 	return packet
	// }

	// rows := ExtractRows(packet)
	return packet
}
