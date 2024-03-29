package main

import (
	"fmt"

	"net"

	"github.com/cakoshakib/distributed-db/storage"
)

/*
GET KEY REQUEST
get [user] [table] [key];

DELETE KEY REQUEST
del [user][table] [key];

ADD KEY, VAL REQUEST
add [user][table] [key] [value];

CREATE TABLE
ct [user] [table];

CREATE USER
cu [user];

RESPONSES
200 OK
201 CREATED (table)
400 BAD REQUEST
404 NOT FOUND (table or key)

userA
- table1
- table2
- table3
userB
- table1
- table2

"userA": {
	"table1": {

	},
	"table2": {

	}
}
*/

func main() {
	fmt.Println("vim-go")

	_, _ = net.Dial("test", "test")

	//storage.AddTable("bob", "table2")
	storage.ReadKV("bob", "table2", "key1")
	storage.AddKV("bob", "table2", "key2", "test")
	storage.ReadKV("bob", "table2", "key2")
	storage.RemoveKV("bob", "table2", "key2")
	storage.ReadKV("bob", "table2", "key2")
}
