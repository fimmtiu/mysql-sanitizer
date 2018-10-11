package main

import (
	"net"

	"github.com/pubnative/mysqlproto-go"
)

// A Client represents a single client connection to the MySQL server.
type ClientConnection struct {
	proxy  *ProxyConnection
	stream *mysqlproto.Stream
}

// NewClientConnection returns a new ClientConnection object.
func NewClientConnection(proxy *ProxyConnection, conn net.Conn) *ClientConnection {
	client := ClientConnection{proxy, nil}
	client.stream = mysqlproto.NewStream(conn)
	return &client
}

// ProcessInput listens for client requests and proxies them to the MySQL server.
func (client *ClientConnection) Run() {
	for {
		packet, err := client.stream.NextPacket()
		if err != nil {
			output.Log("Disconnected from MySQL server: %s", err)
			client.proxy.Close()
			return
		}
		output.Log("Packet: type %d", packet.Payload[0])
		output.Dump(packet.Payload, "Packet contents")
	}
}

func (client *ClientConnection) Close() {
	client.stream.Close()
}
