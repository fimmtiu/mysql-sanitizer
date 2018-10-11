package main

import (
	"net"

	"github.com/pubnative/mysqlproto-go"
)

// A Client represents a single client connection to the MySQL server.
type ClientConnection struct {
	proxy          *ProxyConnection
	stream         *mysqlproto.Stream
	capabilities   uint32
	authPluginData []byte
	database       string
}

type HandshakeContents struct {
	flags          uint32
	maxPacketSize  uint32
	characterSet   byte
	username       string
	password       string
	database       string
	authPluginName string
	connectAttrs   map[string]string
}

// NewClientConnection returns a new ClientConnection object.
func NewClientConnection(proxy *ProxyConnection, conn net.Conn) *ClientConnection {
	var client ClientConnection
	client.proxy = proxy
	client.stream = mysqlproto.NewStream(conn)
	return &client
}

// ProcessInput listens for client requests and proxies them to the MySQL server.
func (client *ClientConnection) Run() {
	var packetCount uint64 = 0
	incoming := make(chan mysqlproto.Packet)
	go client.getPackets(incoming)

	for {
		select {
		case packet := <-client.proxy.ClientChannel:
			packetCount++
			if packetCount == 1 {
				// This is the first packet the server sent, so it must be
				// the start of the handshake.
				client.authPluginData = client.getAuthPluginData(packet)
			}
			WritePacket(client.stream, packet)
		case packet, more := <-incoming:
			if more {
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
			packet = client.replacePassword(packet, config.MysqlUsername, config.MysqlPassword)
		}
		output.Dump(packet.Payload, "Packet %d from client:\n", packetCount)
		channel <- packet
	}
}

func (client *ClientConnection) Close() {
	client.stream.Close()
}

func (client *ClientConnection) replacePassword(packet mysqlproto.Packet, username string, password string) mysqlproto.Packet {
	contents := client.parseHandshakeResponse(packet)
	contents.username = config.MysqlUsername
	contents.password = config.MysqlPassword
	client.database = contents.database

	newPayload := mysqlproto.HandshakeResponse41(
		contents.flags&client.capabilities,
		contents.characterSet,
		contents.username,
		contents.password,
		client.authPluginData,
		contents.database,
		contents.authPluginName,
		map[string]string{}, // FIXME: We don't support client connect attrs yet.
	)
	return mysqlproto.Packet{packet.SequenceID, newPayload[4:]}
}

func (client *ClientConnection) parseHandshakeResponse(packet mysqlproto.Packet) HandshakeContents {
	var contents HandshakeContents
	parser := NewPacketParser(packet)

	contents.flags = parser.ReadFixedInt4()
	contents.maxPacketSize = parser.ReadFixedInt4()
	contents.characterSet = parser.ReadFixedInt1()
	parser.ReadFixedString(23) // ignore reserved null bytes
	contents.username = parser.ReadNullTermString()

	if (contents.flags&mysqlproto.CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA > 0) ||
		(contents.flags&mysqlproto.CLIENT_SECURE_CONNECTION > 0) {
		contents.password = parser.ReadVariableString()
	} else {
		contents.password = parser.ReadNullTermString()
	}

	if contents.flags&mysqlproto.CLIENT_CONNECT_WITH_DB > 0 {
		contents.database = parser.ReadNullTermString()
	}

	if contents.flags&mysqlproto.CLIENT_PLUGIN_AUTH > 0 {
		contents.authPluginName = parser.ReadNullTermString()
	}

	// FIXME: We don't support client connect attrs yet.

	return contents
}

func (client *ClientConnection) getAuthPluginData(packet mysqlproto.Packet) []byte {
	parser := NewPacketParser(packet)
	parser.ReadFixedInt1()                    // protocol version
	parser.ReadNullTermString()               // server version
	parser.ReadFixedInt4()                    // connection id
	data := []byte(parser.ReadFixedString(8)) // initial part of the auth plugin data
	parser.ReadFixedInt1()                    // unused filler byte
	lowerFlags := parser.ReadFixedInt2()      // capability flags

	if uint64(len(packet.Payload)) <= parser.offset {
		return data
	}

	parser.ReadFixedInt1()               // character set
	parser.ReadFixedInt2()               // status flags
	upperFlags := parser.ReadFixedInt2() // more capability flags, sheesh
	client.capabilities = uint32(lowerFlags) | (uint32(upperFlags) << 16)
	var dataLen uint64 = uint64(parser.ReadFixedInt1() - 8)
	if dataLen > 13 {
		dataLen = 13
	}
	parser.ReadFixedString(10) // unused garbage

	if client.capabilities&mysqlproto.CLIENT_SECURE_CONNECTION > 0 {
		// Don't ask about the -1. :~(
		data = append(data, []byte(parser.ReadFixedString(dataLen-1))...)
	}

	return data
}
