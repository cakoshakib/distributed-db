package main

import (
	"context"
	"os"
	"os/signal"

	"go.uber.org/zap"
	//"github.com/cakoshakib/distributed-db/storage"
	"github.com/cakoshakib/distributed-db/network"
	log "github.com/cakoshakib/distributed-db/commons"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger, _ := zap.NewDevelopment()
	ctx = context.WithValue(ctx, log.LoggerKey, logger)
	defer logger.Sync()

	server, err := network.NewServer(ctx)
	if err != nil {
		logger.Error("server: failed initialization with error", zap.Error(err))
		return
	}

	server.Start(ctx)
}
