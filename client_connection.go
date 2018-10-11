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
	incoming := make(chan mysqlproto.Packet)
	go client.getPackets(incoming)

	for {
		select {
		case packet := <-client.proxy.ClientChannel:
			WritePacket(client.stream, packet)
		case packet, more := <-incoming:
			if more {
				output.Dump(packet.Payload, "Packet from client:\n")
				client.proxy.ServerChannel <- packet
			} else {
				client.proxy.Close()
			}
		}
	}
}

func (client *ClientConnection) getPackets(channel chan mysqlproto.Packet) {
	for {
		packet, err := client.stream.NextPacket()
		if err != nil {
			output.Log("Disconnected from client: %s", err)
			close(channel)
			return
		}
		output.Dump(packet.Payload, "Packet from client:\n")
		channel <- packet
	}
}

func (client *ClientConnection) Close() {
	client.stream.Close()
}
