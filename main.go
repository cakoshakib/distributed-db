package main

import (
	"fmt"

	//"github.com/cakoshakib/distributed-db/storage"
	"github.com/cakoshakib/distributed-db/network"
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
	server, err := network.NewServer()
	if err != nil {
		fmt.Println("server failed creation")
		return
	}

	server.Start()
}
