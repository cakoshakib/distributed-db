package storage

import (
	"encoding/json"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

// we need to treat this store as a finite state machine
type Store struct {
	RaftDir  string
	RaftBind string
	raft     *raft.Raft
	logger   *zap.Logger
}

func New(logger *zap.Logger) *Store {
	return &Store{
		logger: logger,
	}
}

func (s *Store) Open(nodeID string) error {
	// open raft store for this node
	return nil
}

func (s *Store) Join(nodeId, addr string) error {
	// add this node to the cluster
	return nil
}

func (s *Store) Apply(log *raft.Log) {
	var req dbrequest.DBRequest
	if err := json.Unmarshal(log.Data, &req); err != nil {
		s.logger.Error("raft: failed to unmarshal request", zap.Error(err))
	}

	switch req.Op {
	case dbrequest.CreateUser:
		if err := s.AddUser(req); err != nil {
			s.logger.Warn("raft.Apply(): unable to add user", zap.String("user", req.User), zap.Error(err))
		}
	case dbrequest.DeleteUser:
		if err := s.DeleteUser(req); err != nil {
			s.logger.Warn("raft.Apply(): unable to delete user", zap.String("user", req.User), zap.Error(err))
		}
	case dbrequest.CreateTable:
		if err := s.AddTable(req); err != nil {
			s.logger.Warn("raft.Apply(): unable to create table", zap.String("table", req.Table), zap.Error(err))
		}
	case dbrequest.DeleteTable:
		if err := s.DeleteTable(req); err != nil {
			s.logger.Warn("raft.Apply(): unable to delete table", zap.String("table", req.Table), zap.Error(err))
		}
	case dbrequest.AddKV:
		if err := s.AddKV(req); err != nil {
			s.logger.Warn("raft.Apply(): unable to create KV", zap.String("key", req.Key), zap.String("value", req.Value), zap.Error(err))
		}
	case dbrequest.GetKV:
		_, err := s.ReadKV(req)
		if err != nil {
			s.logger.Warn("raft.Apply(): unable to get KV", zap.String("key", req.Key), zap.Error(err))
		}
	case dbrequest.DelKV:
		if err := s.RemoveKV(req); err != nil {
			s.logger.Warn("raft.Apply(): unable to delete KV", zap.String("key", req.Key), zap.Error(err))
		}
	default:
		s.logger.Error("raft.Apply(): unrecognized log request", zap.String("operation", string(req.Op)))
		return
	}
}
