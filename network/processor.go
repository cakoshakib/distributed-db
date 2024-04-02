package network

import (
	"bufio"
	"context"
	"net"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/storage"
	"go.uber.org/zap"
)

// BUG: if user sends ; as part of their data, then the reader terminates early with unread content.
// we should standardize something like "\;" to represent a semicolon in user data with escape sequences and modify the logic below accordingly.
func process(ctx context.Context, conn net.Conn) {
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

	req := NewRequest(msg)
	logger.Info(
		"server.process(): received request", zap.String("remoteAddr", remoteAddr),
		zap.String("operation", string(req.op)), zap.String("user", req.user), zap.String("table", req.table), zap.String("key", req.key), zap.String("value", req.value),
	)

	res := handleRequest(ctx, req)
	if _, err := conn.Write([]byte(res + "\n")); err != nil {
		logger.Error("Error writing to connection", zap.String("remoteAddr", remoteAddr), zap.String("response", res), zap.Error(err))
	}
}

// TODO: be able to pass in DBRequest right into storage functions, such as storage.CreateUser(req)
// TODO: perhaps make response enums to avoid magic strings
func handleRequest(ctx context.Context, req DBRequest) string {
	logger := log.LoggerFromContext(ctx)

	if !req.Validate() {
		logger.Info(
			"server.handleRequest(): bad request was received",
			zap.String("operation", string(req.op)), zap.String("user", req.user), zap.String("table", req.table), zap.String("key", req.key), zap.String("value", req.value),
		)
		return "400 BAD REQUEST"
	}

	switch req.op {
	case CreateUser:
		if err := storage.AddUser(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to add user", zap.String("user", req.user), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case DeleteUser:
		if err := storage.DeleteUser(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to delete user", zap.String("user", req.user), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case CreateTable:
		if err := storage.CreateTable(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to create table", zap.String("table", req.table), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case DeleteTable:
		if err := storage.DeleteTable(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to delete table", zap.String("table", req.table), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case AddKV:
		if err := storage.AddKV(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to create KV", zap.String("key", req.key), zap.String("value", req.value) zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case GetKV:
		if err := storage.GetKV(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to get KV", zap.String("key", req.key), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	case DelKV:
		if err := storage.DelKV(ctx, req); err != nil {
			logger.Warn("server.handleRequest(): unable to delete KV", zap.String("key", req.key), zap.Error(err))
			return "401 REQUEST FAILED"
		}
		return "200 OK"
	default:
		return "400 BAD REQUEST"
	}
}
