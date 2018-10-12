package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"github.com/pubnative/mysqlproto-go"
)

// ServerConnection is a connection to the MySQL server.
type ServerConnection struct {
	proxy      *ProxyConnection
	stream     *mysqlproto.Stream
	sanitizing bool
	finished   bool
}

// NewServerConnection returns a ServerConnection that's connected to the MySQL server.
func NewServerConnection(proxy *ProxyConnection) (*ServerConnection, error) {
	server := ServerConnection{proxy, nil, false, false}

	addrString := config.MysqlHost + ":" + strconv.Itoa(config.MysqlPort)
	addr, err := net.ResolveTCPAddr("tcp", addrString)
	if err != nil {
		return nil, fmt.Errorf("Can't resolve host %s: %s", config.MysqlHost, err)
	}
	addr.Port = config.MysqlPort

	socket, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to %s on port %d:  %s", config.MysqlHost, addr.Port, err)
	}
	server.stream = mysqlproto.NewStream(socket)

	return &server, nil
}

func (server *ServerConnection) ToggleSanitizing(active bool) {
	server.sanitizing = active
}

func (server *ServerConnection) Run() {
	defer server.proxy.Close()
	server.doHandshake()

	for !server.finished {
		packet := <-server.proxy.ServerChannel
		WritePacket(server.stream, packet)

		if packetCommand(packet) == mysqlproto.COM_QUERY {
			server.handleQueryResponse()
		} else {
			server.handleOtherResponse()
		}
	}
}

// Close closes the connection to the MySQL server.
func (server *ServerConnection) Close() {
	server.stream.Close()
}

func (server *ServerConnection) doHandshake() {
	welcomePacket, err := server.stream.NextPacket()
	output.Dump(welcomePacket.Payload, "Welcome packet from server:\n")
	if err != nil {
		output.Log("Couldn't complete handshake to MySQL server: %s", err)
		server.finished = true
		return
	}
	server.proxy.ClientChannel <- welcomePacket

	clientHandshake := <-server.proxy.ServerChannel
	WritePacket(server.stream, clientHandshake)

	response, err := server.stream.NextPacket()
	output.Dump(response.Payload, "Handshake response packet from server:\n")

	if err != nil {
		output.Log("Couldn't complete handshake to MySQL server: %s", err)
		server.finished = true
		return
	}
	if !packetIsOK(response) {
		output.Log("Bad handshake response from MySQL server")
		server.finished = true
		return
	}
	server.proxy.ClientChannel <- response
}

func (server *ServerConnection) handleQueryResponse() {
	for {
		response, err := server.stream.NextPacket()
		if err != nil {
			output.Log("Couldn't receive packet from MySQL server: %s", err)
			server.finished = true
			return
		}
		output.Dump(response.Payload, "Packet from server:\n")

		if packetIsOK(response) || packetIsERR(response) || packetIsEOF(response) {
			server.proxy.ClientChannel <- response
			break
		} else {
			columns, err := server.readColumnDefinitions(response)
			if err != nil {
				output.Log("Couldn't receive column definitions from MySQL server: %s", err)
				server.finished = true
				return
			}

			eofPacket, err := server.stream.NextPacket()
			if err != nil {
				output.Log("Couldn't receive column definitions from MySQL server: %s", err)
				server.finished = true
				return
			}
			output.Dump(eofPacket.Payload, "End of column definitions packet from server:\n")
			server.proxy.ClientChannel <- eofPacket

			for {
				rowPacket, err := server.stream.NextPacket()
				output.Dump(rowPacket.Payload, "Response packet from server:\n")

				if err != nil {
					output.Log("Couldn't receive column definitions from MySQL server: %s", err)
					server.finished = true
					return
				}
				if packetIsOK(rowPacket) || packetIsERR(rowPacket) || packetIsEOF(rowPacket) {
					server.proxy.ClientChannel <- rowPacket
					return
				}

				rows, err := readRowValues(rowPacket, columns)
				if err != nil {
					output.Log("Couldn't receive row values from MySQL server: %s", err)
					server.finished = true
					return
				}

				server.proxy.ClientChannel <- constructNewResponse(rowPacket, rows)
			}
		}
	}
}

func (server *ServerConnection) handleOtherResponse() {
	for {
		response, err := server.stream.NextPacket()
		if err != nil {
			output.Log("Couldn't receive packet from MySQL server: %s", err)
			server.finished = true
			return
		}
		server.proxy.ClientChannel <- response
		if packetIsOK(response) || packetIsERR(response) || packetIsEOF(response) {
			break
		}
	}
}

func packetIsOK(packet mysqlproto.Packet) bool {
	return packet.Payload[0] == 0 && len(packet.Payload) >= 7
}

func packetIsERR(packet mysqlproto.Packet) bool {
	return packet.Payload[0] == 0xFF
}

func packetIsEOF(packet mysqlproto.Packet) bool {
	return packet.Payload[0] == 0xFE && len(packet.Payload) < 9
}

func packetCommand(packet mysqlproto.Packet) byte {
	return packet.Payload[0]
}

func (server *ServerConnection) readColumnDefinitions(packet mysqlproto.Packet) ([]Column, error) {
	parser := NewPacketParser(packet)
	columnCount := parser.ReadEncodedInt()

	columns := make([]Column, columnCount)
	server.proxy.ClientChannel <- packet

	for i := 0; i < int(columnCount); i++ {
		packet, err := server.stream.NextPacket()
		if err != nil {
			return nil, err
		}
		output.Dump(packet.Payload, "Column definition packet from server:\n")
		parser = NewPacketParser(packet)
		server.proxy.ClientChannel <- packet

		column, err := ReadColumn(parser)
		if err != nil {
			return nil, err
		}
		columns[i] = column
	}
	return columns, nil
}

func readRowValues(packet mysqlproto.Packet, columns []Column) ([][]byte, error) {
	parser := NewPacketParser(packet)
	rows := [][]byte{}

	for _, col := range columns {
		value, nonNull := parser.ReadStringOrNull()
		if nonNull {
			rowVal := []byte(value)
			if !col.IsSafe() {
				rowVal = sanitizeRow(rowVal, col)
			}
			rows = append(rows, rowVal)
		} else {
			rows = append(rows, nil)
		}
	}

	return rows, nil
}

func sanitizeRow(row []byte, column Column) []byte {
	sum := sha256.Sum256(row)
	newRow := make([]byte, sha256.Size*2)
	hex.Encode(newRow, sum[:])

	if uint32(len(newRow)) > column.Length {
		newRow = newRow[:column.Length]
	}
	return newRow
}

func constructNewResponse(originalPacket mysqlproto.Packet, rows [][]byte) mysqlproto.Packet {
	newPacket := mysqlproto.Packet{originalPacket.SequenceID, []byte{}}

	for _, row := range rows {
		row = append(LengthEncodedInt(uint(len(row))), row...)
		newPacket.Payload = append(newPacket.Payload, row...)
	}

	return newPacket
}
