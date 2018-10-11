package main

import (
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
}

// NewServerConnection returns a ServerConnection that's connected to the MySQL server.
func NewServerConnection(proxy *ProxyConnection) (*ServerConnection, error) {
	server := ServerConnection{proxy, nil, false}

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
	incoming := make(chan mysqlproto.Packet)
	go server.getPackets(incoming)

	for {
		select {
		case packet := <-server.proxy.ServerChannel:
			WritePacket(server.stream, packet)
		case packet, more := <-incoming:
			if more {
				server.proxy.ClientChannel <- packet
			} else {
				server.proxy.Close()
			}
		}
	}
}

func (server *ServerConnection) getPackets(channel chan mysqlproto.Packet) {
	for {
		packet, err := server.stream.NextPacket()
		if err != nil {
			output.Log("Disconnected from MySQL server: %s", err)
			close(channel)
			return
		}
		output.Dump(packet.Payload, "Packet from server:\n")
		channel <- packet
	}
}

// Close closes the connection to the MySQL server.
func (server *ServerConnection) Close() {
	server.stream.Close()
}
