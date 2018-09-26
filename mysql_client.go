package main

import (
	"fmt"
	"net"

	"github.com/pubnative/mysqlproto-go"
)

// MysqlClient is a connection to the MySQL server.
type MysqlClient struct {
	// conn mysqlproto.Conn
	stream *mysqlproto.Stream
}

// NewMysqlClient returns a MysqlClient that's connected to the MySQL server.
func NewMysqlClient(config Config, capabilityFlags uint32, database string) (*MysqlClient, error) {
	var mc MysqlClient

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

	// mc.conn, err = mysqlproto.ConnectPlainHandshake(socket, capabilityFlags, config.MysqlUsername, config.MysqlPassword, database, map[string]string{})
	// if err != nil {
	// 	return nil, fmt.Errorf("Handshake to MySQL server failed: %s", err)
	// }

	return &mc, nil
}

// Close closes the connection to the MySQL server.
func (mc *MysqlClient) Close() {
	mc.stream.Close()
}
