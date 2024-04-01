package network

import (
	"fmt"
	"net"
)

func process(conn net.Conn) {
	fmt.Printf("server.process(): received connection from %s", conn.RemoteAddr())

	// TODO

	fmt.Printf("server.process(): closing connection with %s", conn.RemoteAddr())
	conn.Close()
}
