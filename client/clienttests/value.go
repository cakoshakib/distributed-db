package clienttests

import (
	"fmt"
	"sync"
	"time"
)

var userTables = map[string]map[string]map[string]string{}
var mutex = &sync.RWMutex{}

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
			return "\"" + value + "\"\n", true
		}
	}
	return "", false
}

func ValueTest() {
	total := 0
	miss := 0
	user := "user1"
	table := "table1"

	add_user(user)
	add_table(user, table)
	ProcessRequest(fmt.Sprintf("cu %s;\n", user), leader)
	ProcessRequest(fmt.Sprintf("ct %s %s;\n", user, table), leader)

	n := 1000

	for i := 1; i <= n; i++ {
		add_kv(user, table, fmt.Sprintf("test%d", i), fmt.Sprintf("value%d", i))
		request := fmt.Sprintf("add %s %s test%d value%d;\n", user, table, i, i)
		ProcessRequest(request, leader)

		request = fmt.Sprintf("get %s %s test%d;\n", user, table, i)
		localVal, _ := check_kv(user, table, fmt.Sprintf("test%d", i))
		time.Sleep(100 * time.Millisecond)
		serverVal, _ := ProcessRequest(request, follower)

		if serverVal != localVal {
			fmt.Println(fmt.Sprintf("local: %s, server %s", localVal, serverVal))
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
