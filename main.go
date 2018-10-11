package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pubnative/mysqlproto-go"
)

var output Output
var config Config

func main() {
	config = GetConfig()
	output = NewOutput(config)
	listener := openListeningSocket(config.ListeningPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Can't accept incoming connection on port %d: %s", config.ListeningPort, err)
		}

		proxy, err := NewProxyConnection(conn)
		if err == nil {
			proxy.Start()
		} else {
			output.Log("Can't open connection to %s: %s", config.MysqlHost, err)
		}
	}
}

// Returns a TCP socket that's listening on the given port.
func openListeningSocket(port int) net.Listener {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("Can't listen on port %d: %s", port, err)
	}
	return listener
}

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
