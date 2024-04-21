package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	user := "user1"
	table := "table1"

	processRequest(fmt.Sprintf("cu %s;\n", user))
	processRequest(fmt.Sprintf("ct %s %s;\n", user, table))

	//Begin test
	n := 1000

	start := time.Now()
	for i := 1; i <= n; i++ {
		request := fmt.Sprintf("add %s %s test%d value%d;\n", user, table, i, i)
		processRequest(request)
	}
	checkpoint := time.Now()
	for i := 1; i <= n; i++ {
		request := fmt.Sprintf("get %s %s test%d;\n", user, table, i)
		processRequest(request)
	}
	end := time.Now()

	elapsed_write := checkpoint.Sub(start)
	elapsed_read := end.Sub(checkpoint)
	fmt.Println("Time taken for writes:", elapsed_write)
	fmt.Println("Time taken for reads:", elapsed_read)
}

func processRequest(req string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
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
