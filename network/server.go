package network

import (
	"context"
	"errors"
	"net"

	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/commons/metricsink"
	"github.com/cakoshakib/distributed-db/storage"
	"go.uber.org/zap"
)

type server struct {
	listener      net.Listener
	store         *storage.Store
	metrics       *metricsink.MetricHandler
	cancelMetrics context.CancelFunc
}

func NewServer(ctx context.Context, port string, store *storage.Store, metricPath string, doMetrics bool) (server, error) {
	server := server{}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return server, err
	}
	server.listener = listener

	server.store = store

	// init metrics sink
	if doMetrics {
		metricHandler, err := metricsink.NewMetricHandler(metricPath)
		if err != nil {
			return server, err
		}
		server.metrics = metricHandler
	}

	return server, nil
}

func (s server) Start(ctx context.Context) {
	logger := log.LoggerFromContext(ctx)
	logger.Info("server.start(): Starting server")

	// s.metrics != nil checks if we should be tracking metrics
	// aka, we metrics are set to true
	if s.metrics != nil {
		logger.Info("metric handler active")
		metricsCtx, cancel := context.WithCancel(context.Background())
		metricsCtx = context.WithValue(metricsCtx, log.LoggerKey, logger)
		s.cancelMetrics = cancel
		go s.metrics.LogMetricsToCSV(metricsCtx)
	}

	go func() {
		<-ctx.Done()
		logger.Info("Context is cancelled; Stopping server")
		s.Stop(ctx)
		if s.cancelMetrics != nil {
			s.cancelMetrics()
		}
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
		go process(ctx, conn, s.store)
	}
}

func (s server) Stop(ctx context.Context) {
	logger := log.LoggerFromContext(ctx)
	logger.Info("server.close(): closing server")
	if err := s.listener.Close(); err != nil {
		logger.Error("server.close(): error closing server", zap.Error(err))
	}
}
