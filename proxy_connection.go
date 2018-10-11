package main

import (
	"net"

	"github.com/pubnative/mysqlproto-go"
)

type ProxyConnection struct {
	client        *ClientConnection
	server        *ServerConnection
	ClientChannel chan mysqlproto.Packet
	ServerChannel chan mysqlproto.Packet
	Capabilities  uint32
	Database      string
}

func NewProxyConnection(conn net.Conn) (*ProxyConnection, error) {
	var err error
	var proxy ProxyConnection
	proxy.ClientChannel = make(chan mysqlproto.Packet)
	proxy.ServerChannel = make(chan mysqlproto.Packet)

	proxy.client = NewClientConnection(&proxy, conn)
	proxy.server, err = NewServerConnection(&proxy)
	if err != nil {
		return nil, err
	}

	return &proxy, nil
}

func (proxy *ProxyConnection) Start() {
	go proxy.client.Run()
	go proxy.server.Run()
}

func (proxy *ProxyConnection) Close() {
	proxy.client.Close()
	proxy.server.Close()
}
