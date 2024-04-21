package clienttests

import (
	"bufio"
	"fmt"
	"net"
)

const (
	leader   = "8080"
	follower = "8070"
)

func ProcessRequest(req string, port string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Handle Error Later")
		return "", err
	}

	conn.Write([]byte((req)))
	reader := bufio.NewReader(conn)

	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("server.process(): error reading from connection")
		return "", err
	}

	return msg, nil
}
