package network

import (
	"bufio"
	"context"
	"net"

	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"github.com/cakoshakib/distributed-db/storage"
	"go.uber.org/zap"
)

// BUG: if user sends ; as part of their data, then the reader terminates early with unread content.
// we should standardize something like "\;" to represent a semicolon in user data with escape sequences and modify the logic below accordingly.
func process(ctx context.Context, conn net.Conn, store storage.Store) {
	logger := log.LoggerFromContext(ctx)
	remoteAddr := conn.RemoteAddr().String()
	logger.Info("server.process(): received connection", zap.String("remoteAddr", remoteAddr))
	defer conn.Close()
	defer logger.Info("server.process(): closing connection", zap.String("remoteAddr", remoteAddr))

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString(';')
	if err != nil {
		logger.Error("server.process(): error reading from connection", zap.Error(err))
		return
	}

	req := dbrequest.NewRequest(msg)
	logger.Info(
		"server.process(): received request", zap.String("remoteAddr", remoteAddr),
		zap.String("operation", string(req.Op)), zap.String("user", req.User), zap.String("table", req.Table), zap.String("key", req.Key), zap.String("value", req.Value),
	)

	res := handleRequest(ctx, req, store)
	if _, err := conn.Write([]byte(res + "\n")); err != nil {
		logger.Error("Error writing to connection", zap.String("remoteAddr", remoteAddr), zap.String("response", res), zap.Error(err))
	}
}

// TODO: be able to pass in DBRequest right into storage functions, such as storage.CreateUser(req)
// TODO: perhaps make response enums to avoid magic strings
func handleRequest(ctx context.Context, req dbrequest.DBRequest, store storage.Store) string {
	logger := log.LoggerFromContext(ctx)

	if !req.Validate() {
		logger.Info(
			"server.handleRequest(): bad request was received",
			zap.String("operation", string(req.Op)), zap.String("user", req.User), zap.String("table", req.Table), zap.String("key", req.Key), zap.String("value", req.Value),
		)
		return "400 BAD REQUEST"
	}

	err := store.HandleRequest(req)
	// TODO: should check what type of error
	if err != nil {
		logger.Warn("server.handleRequest(): error handling request", zap.String("op", string(req.Op)))
		return "401 REQUEST FAILED"
	}
	return "200 OK"

}
