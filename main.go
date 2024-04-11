package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"
	//"github.com/cakoshakib/distributed-db/storage"
	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/network"
)

const (
	DefaultRaftAddr = "localhost:12000"
	DefaultTCPPort  = "8080"
)

var (
	tcpPort  string
	nodeID   string
	raftAddr string
	joinAddr string
)

func init() {
	flag.StringVar(&tcpPort, "tcpPort", DefaultTCPPort, "TCP Port for client requests")
	flag.StringVar(&nodeID, "id", "", "Raft Node ID")
	flag.StringVar(&raftAddr, "raftAddr", DefaultRaftAddr, "Raft binding address")
	flag.StringVar(&joinAddr, "joinAddr", "", "Client-facing address to join Raft cluster")
}

func main() {
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger, _ := zap.NewDevelopment()
	ctx = context.WithValue(ctx, log.LoggerKey, logger)
	defer logger.Sync()
	logger.Info(fmt.Sprintf("received params %s %s %s %s", tcpPort, nodeID, raftAddr, joinAddr))

	server, err := network.NewServer(ctx, tcpPort)
	if err != nil {
		logger.Error("server: failed initialization with error", zap.Error(err))
		return
	}

	server.Start(ctx)
}
