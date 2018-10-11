package main

import (
	"fmt"
	"net"

	"github.com/pubnative/mysqlproto-go"
)

// MysqlClient is a connection to the MySQL server.
type MysqlClient struct {
	Input      chan mysqlproto.Packet
	Done       chan error
	stream     *mysqlproto.Stream
	sanitizing bool
}

// NewMysqlClient returns a MysqlClient that's connected to the MySQL server.
func NewMysqlClient(config Config) (*MysqlClient, error) {
	mc := MysqlClient{nil, false}

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

	mc.Input = make(chan mysqlproto.Packet)
	mc.Done = make(chan error)

	return &mc, nil
}

func (mc *MysqlClient) ToggleSanitizing(active bool) {
	mc.sanitizing = active
}

func (mc *MysqlClient) Run() {
	// ...
}

// Close closes the connection to the MySQL server.
func (mc *MysqlClient) Close() {
	mc.stream.Close()
}
