package main

import (
	"fmt"
	"log"
	"net"
)

var output Output

func main() {
	config := GetConfig()
	output = NewOutput(config)
	listener := openListeningSocket(config.ListeningPort)

	/* We want:
	  - A sanitizer function which takes a packet and returns a sanitized packet
	  - N connections to the MySQL server
		- when instantiated, creates a goroutine and returns a channel to communicate with it
		- listens to the channel for packets from the client
		- when it gets one, sends it to the MySQL server
		- listens for the response
		- runs each response packet through the sanitizer
		- sends each sanitized response packet back over the channel
	  - N client connections
		- use mysqlproto to read incoming packets from the connection
		- sends the packet to its MySQL server connection channel
	*/

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Can't accept incoming connection on port %d: %s", config.ListeningPort, err)
		}

		client := NewClient(config, conn)
		go client.ProcessInput()
	}
}

// Returns a TCP socket that's listening on the given port.
func openListeningSocket(port int) net.Listener {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("Can't listen on port %d: %s", port, err)
	}
	return listener
}
