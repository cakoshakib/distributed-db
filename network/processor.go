package network

import (
	"bufio"
	"context"
	"fmt"
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

	switch req.Op {
	case dbrequest.CreateUser:
		if err := store.AddUser(req); err != nil {
			logger.Warn("server.handleRequest(): unable to add user", zap.String("user", req.User), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case dbrequest.DeleteUser:
		if err := store.DeleteUser(req); err != nil {
			logger.Warn("server.handleRequest(): unable to delete user", zap.String("user", req.User), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case dbrequest.CreateTable:
		if err := store.AddTable(req); err != nil {
			logger.Warn("server.handleRequest(): unable to create table", zap.String("table", req.Table), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case dbrequest.DeleteTable:
		if err := store.DeleteTable(req); err != nil {
			logger.Warn("server.handleRequest(): unable to delete table", zap.String("table", req.Table), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case dbrequest.AddKV:
		if err := store.AddKV(req); err != nil {
			logger.Warn("server.handleRequest(): unable to create KV", zap.String("key", req.Key), zap.String("value", req.Value), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case dbrequest.GetKV:
		val, err := store.ReadKV(req)
		if err != nil {
			logger.Warn("server.handleRequest(): unable to get KV", zap.String("key", req.Key), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return fmt.Sprintf("%s\n200 OK", val)
	case dbrequest.DelKV:
		if err := store.RemoveKV(req); err != nil {
			logger.Warn("server.handleRequest(): unable to delete KV", zap.String("key", req.Key), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	default:
		return "400 BAD REQUEST"
	}
}
