package main

import (
	"fmt"
    "sync"
)

var userTables = map[string]map[string]map[string]string{}
var mutex = &sync.RWMutex{}
var count int = 0
var misses int = 0

func add_user(user string){
	mutex.Lock()
	userTables[user] = make(map[string]map[string]string)
	mutex.Unlock()
}

func add_table(user string, table string){
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

func check_kv(user string, table string, key string)(string, bool){
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
	add_user("nick")
	add_table("nick", "table1")
	add_kv("nick", "table1", "key1", "urmom")
	value, exists := check_kv("nick", "table1", "key1")
	if exists {
		fmt.Println(value) // prints: value1
	}
}