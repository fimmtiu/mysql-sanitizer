package main

import (
	"fmt"
	"net"

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
	mc := ServerConnection{proxy, nil, false}

	addr, err := net.ResolveTCPAddr("tcp", config.MysqlHost)
	if err != nil {
		return nil, fmt.Errorf("Can't resolve host %s: %s", config.MysqlHost, err)
	}
	addr.Port = config.MysqlPort

	socket, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to %s on port %d:  %s", config.MysqlHost, addr.Port, err)
	}
	mc.stream = mysqlproto.NewStream(socket)

	return &mc, nil
}

func (mc *ServerConnection) ToggleSanitizing(active bool) {
	mc.sanitizing = active
}

func (mc *ServerConnection) Run() {
	// ...
}

// Close closes the connection to the MySQL server.
func (mc *ServerConnection) Close() {
	mc.stream.Close()
}
