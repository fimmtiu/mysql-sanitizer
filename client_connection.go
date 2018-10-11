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

type HandshakeContents struct {
	flags          uint32
	characterSet   byte
	username       string
	password       string
	authPluginData []byte
	database       string
	authPluginName string
	connectAttrs   map[string]string
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
				return
			}
		}
	}
}

func (client *ClientConnection) getPackets(channel chan mysqlproto.Packet) {
	var packetCount uint64 = 0

	for {
		packet, err := client.stream.NextPacket()
		if err != nil {
			output.Log("Disconnected from client: %s", err)
			close(channel)
			return
		}
		packetCount++
		if packetCount == 1 {
			// This is the first packet the client sent, so it must be a handshake.
			client.replacePassword(&packet, config.MysqlUsername, config.MysqlPassword)
		}
		output.Dump(packet.Payload, "Packet from client:\n")
		channel <- packet
	}
}

func (client *ClientConnection) Close() {
	client.stream.Close()
}

func (client *ClientConnection) replacePassword(packet *mysqlproto.Packet, username string, password string) {
	// contents := client.parseHandshakeResponse(packet)

	// packet.Payload = contents
}

func (client *ClientConnection) parseHandshakeReponse(packet *mysqlproto.Packet) HandshakeContents {
	var contents HandshakeContents
	parser := NewPacketParser(packet)
	contents.flags = parser.ReadFixedInt4()
	contents.characterSet = parser.ReadFixedInt1()
	// contents.username

	return contents
}
