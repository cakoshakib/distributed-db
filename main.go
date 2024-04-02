package main

import (
	"context"
	"os"
	"os/signal"

	"go.uber.org/zap"
	//"github.com/cakoshakib/distributed-db/storage"
	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/network"
)

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
