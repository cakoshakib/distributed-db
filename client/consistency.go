package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Print("Handle Error Later")
		return
	}
	defer conn.Close()

	request := "cu nick;"

	fmt.Fprintf(conn, request)

	
}
