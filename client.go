package main

import (
	"net"

	"github.com/pubnative/mysqlproto-go"
)

// A Client represents a single client connection to the MySQL server.
type Client struct {
	stream *mysqlproto.Stream
	mysql  *MysqlClient
}

// NewClient returns a new Client object.
func NewClient(config Config, conn net.Conn) (*Client, error) {
	var client Client
	var err error

	client.stream = mysqlproto.NewStream(conn)
	client.mysql, err = NewMysqlClient(config)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// ProcessInput listens for client requests and proxies them to the MySQL server.
func (client *Client) ProcessInput() {
	for {
		packet, err := client.stream.NextPacket()
		if err != nil {
			output.Log("Disconnected from MySQL server: %s", err)
			client.close()
			return
		}
		output.Log("Packet: type %d", packet.Payload[0])
		output.Dump(packet.Payload, "Packet contents")
	}
}

func (client *Client) close() {
	client.stream.Close()
}
