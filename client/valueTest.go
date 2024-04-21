package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var userTables = map[string]map[string]map[string]string{}
var mutex = &sync.RWMutex{}
var count int = 0
var misses int = 0

func add_user(user string) {
	mutex.Lock()
	userTables[user] = make(map[string]map[string]string)
	mutex.Unlock()
}

func add_table(user string, table string) {
	mutex.Lock()
	if _, ok := userTables[user]; ok {
		userTables[user][table] = make(map[string]string)
	}
	mutex.Unlock()
}

func add_kv(user string, table string, key string, value string) {
	mutex.Lock()
	if _, ok := userTables[user]; ok {
		if _, ok := userTables[user][table]; ok {
			userTables[user][table][key] = value
		}
	}
	mutex.Unlock()
}

func check_kv(user string, table string, key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	if _, ok := userTables[user]; ok {
		if value, ok := userTables[user][table][key]; ok {
			return value, true
		}
	}
	return "", false
}

func main() {
	total := 0
	miss := 0

	value, exists := check_kv("nick", "table1", "key1")
	if exists {
		fmt.Println(value) // prints: value1
	}

	user := "user1"
	table := "table1"
	add_user(user)
	add_table(user, table)
	processRequest(fmt.Sprintf("cu %s;\n", user))
	processRequest(fmt.Sprintf("ct %s %s;\n", user, table))

	n := 1000

	for i := 1; i <= n; i++ {
		request := fmt.Sprintf("add %s %s test%d value%d;\n", user, table, i, i)
		processRequest(request)

		request = fmt.Sprintf("get %s %s test%d;\n", user, table, i)
		serverVal, _ := processRequest(request)
		localVal, exists := check_kv(user, table, fmt.Sprintf("test%d", i))

		if !exists || serverVal != localVal {
			miss++
		}
		total++
	}

	missRate := float64(miss) / float64(total)
	fmt.Println("Hits: ", total-miss)
	fmt.Println("Misses: ", miss)
	fmt.Println("Total: ", total)

	fmt.Println("Miss Rate: ", missRate)
	fmt.Println("Hit Rate: ", 1-missRate)

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
