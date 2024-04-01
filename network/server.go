package network

import (
	"context"
	"errors"
	"net"
	"os"
	
	"go.uber.org/zap"
	log "github.com/cakoshakib/distributed-db/commons"
)

type server struct {
	listener net.Listener
	//logger   interface{}
}

func NewServer(ctx context.Context) (server, error) {
	server := server{}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return server, err
	}
	server.listener = listener

	return server, nil
}

func (s server) Start(ctx context.Context) {
	logger := log.LoggerFromContext(ctx)
	logger.Info("server.start(): Starting server")

	go func() {
		<-ctx.Done()
		logger.Info("Context is cancelled; Stopping server")
		s.Stop(ctx)
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Info("server.start(): listener closed")
				break
			}
			
			logger.Error("server.start() error", zap.Error(err))
		}

		go process(ctx, conn)
	}
}

func (s server) Stop(ctx context.Context) {
	logger := log.LoggerFromContext(ctx)
	logger.Info("server.close(): closing server")
	if err := s.listener.Close(); err != nil {
		logger.Error("server.close(): error closing server", zap.Error(err))
		os.Exit(1)
	}
}
