package network

import (
	"context"
	"net"

	"go.uber.org/zap"
	log "github.com/cakoshakib/distributed-db/commons"
)

func process(ctx context.Context, conn net.Conn) {
	logger := log.LoggerFromContext(ctx)
	logger.Info("server.process(): received connection", zap.String("remoteAddr", conn.RemoteAddr().String()))

	// TODO

	logger.Info("server.process(): closing connection", zap.String("remoteAddr", conn.RemoteAddr().String()))
	conn.Close()
}
