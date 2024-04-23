package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"

	"go.uber.org/zap"
	//"github.com/cakoshakib/distributed-db/storage"
	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/network"
	"github.com/cakoshakib/distributed-db/storage"
	bdbclient "github.com/cakoshakib/distributed-db/storage/boltdbclient"
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
	dataDir  string
)

func init() {
	flag.StringVar(&tcpPort, "tcpPort", DefaultTCPPort, "TCP Port for client requests")
	flag.StringVar(&raftAddr, "raftAddr", DefaultRaftAddr, "Raft binding address")
	flag.StringVar(&nodeID, "id", raftAddr, "Raft Node ID (raftAddr by default)")
	flag.StringVar(&joinAddr, "joinAddr", "", "Client-facing address to join Raft cluster")
	flag.StringVar(&dataDir, "dataDir", "data/", "Directory to store data")
}

func main() {
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// init logger
	logger, _ := zap.NewDevelopment()
	ctx = context.WithValue(ctx, log.LoggerKey, logger)
	defer logger.Sync()
	logger.Info("received params",
		zap.String("tcpPort", tcpPort), zap.String("nodeID", nodeID), zap.String("raftAddr", raftAddr), zap.String("joinAddr", joinAddr), zap.String("dataDir", dataDir),
	)

	// init BoltDB
	db, err := bdbclient.NewBoltDB(path.Join(dataDir, nodeID+".db"))
	if err != nil {
		logger.Fatal("boltdb errors", zap.Error(err))
	}
	defer db.Close()

	// init Raft store
	store := storage.New(logger, dataDir)
	store.RaftBind = raftAddr
	if err := store.Open(joinAddr == "", nodeID, db); err != nil {
		logger.Error("raft: failed to open store", zap.Error(err))
	}

	// init client-facing server
	server, err := network.NewServer(ctx, tcpPort, store)
	if err != nil {
		logger.Error("server: failed initialization with error", zap.Error(err))
		return
	}

	// send join request if necessary
	if joinAddr != "" {
		if err := joinCluster(logger); err != nil {
			logger.Fatal("raft: failed to join cluster", zap.String("joinAddr", joinAddr), zap.Error(err))
		}
	}

	server.Start(ctx)
}

func joinCluster(logger *zap.Logger) error {
	// dial the join address
	conn, err := net.Dial("tcp", joinAddr)
	if err != nil {
		logger.Error("joinCluster(): could not dial join address", zap.String("joinAddr", joinAddr))
		return err
	}
	defer conn.Close()
	payload := []byte(fmt.Sprintf("join %s %s;", nodeID, raftAddr))
	if _, err := conn.Write(payload); err != nil {
		logger.Error("joinCluster(): error sending join request", zap.String("joinAddr", joinAddr), zap.Error(err))
		return err
	}
	return nil
}
