package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

const raftTimeout = 10 * time.Second

// we need to treat this store as a finite state machine
type Store struct {
	RaftBind string
	dataDir  string
	raft     *raft.Raft
	logger   *zap.Logger
}

func New(logger *zap.Logger, dataDir string) *Store {
	return &Store{
		logger:  logger,
		dataDir: dataDir,
	}
}

func (s *Store) Open(firstNode bool, nodeID string) error {
	// open raft store for this node
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// probably want to persist these stores
	snapshots := raft.NewInmemSnapshotStore()
	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	ra, err := raft.NewRaft(config, s, logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	s.raft = ra

	if firstNode {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		s.logger.Info("raft.Open(): bootstrapping cluster", zap.String("id", string(config.LocalID)), zap.String("address", string(transport.LocalAddr())), zap.Any("config", configuration))
		s.raft.BootstrapCluster(configuration)
	}

	return nil
}

func (s *Store) Join(nodeID, addr string) error {
	// add this node to the cluster
	s.logger.Info("raft.Join(): received join request", zap.String("nodeId", nodeID), zap.String("address", addr))
	config := s.raft.GetConfiguration()
	if err := config.Error(); err != nil {
		s.logger.Error("raft.Join(): failed to get config", zap.Error(err))
		return err
	}
	s.logger.Info("raft.Join(): got config")

	raftID := raft.ServerID(nodeID)
	raftAddr := raft.ServerAddress(addr)

	// remove node with given ID if its there
	for _, srv := range config.Configuration().Servers {
		if srv.ID == raftID || srv.Address == raftAddr {
			if srv.ID == raft.ServerID(nodeID) && srv.Address == raftAddr {
				s.logger.Info("raft.Join(): node is already part of cluster, ignoring join", zap.String("nodeId", nodeID), zap.String("address", addr))
				return nil
			}
			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	indexFuture := s.raft.AddVoter(raftID, raftAddr, 0, 0)
	if err := indexFuture.Error(); err != nil {
		s.logger.Error("raft.Join(): failed to add voter", zap.Error(err))
		return err
	}

	s.logger.Info("raft.Join(): successful join", zap.String("nodeId", nodeID), zap.String("address", addr))
	return nil
}

func (s *Store) Restore(rc io.ReadCloser) error {
	// restores store from clean state
	// clean data folder
	dataDir, err := ioutil.ReadDir(s.dataDir)
	if err != nil {
		s.logger.Error("raft.Restore(): could not read data dir", zap.Error(err))
	}
	for _, d := range dataDir {
		os.RemoveAll(user_path(s.dataDir, d.Name()))
	}
	// run through each DBRequest and apply
	decoder := json.NewDecoder(rc)
	for decoder.More() {
		var op dbrequest.DBRequest
		err := decoder.Decode(&op)
		if err != nil {
			return fmt.Errorf("could not decode restore to op: %s", err)
		}
		// not exactly sure if applying logs is a valid way to restore state but let's hope it is :3
		s.HandleRequest(op)
	}

	return rc.Close()
}

func (s *Store) HandleRequest(req dbrequest.DBRequest) error {
	// only leader can take write requests
	s.logger.Info("Handling request", zap.String("op", string(req.Op)), zap.Bool("isWrite", req.IsWrite))
	if req.IsWrite && s.raft.State() != raft.Leader {
		return fmt.Errorf("write op not sent to leader")
	}
	// no need to write reads to log
	if !req.IsWrite {
		_, err := s.ReadKV(req)
		if err != nil {
			s.logger.Warn("raft.HandleRequest(): unable to get KV", zap.String("key", req.Key), zap.Error(err))
			return err
		}
		return nil
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// applies Raft log entry (dbrequest) to store
func (s *Store) Apply(log *raft.Log) interface{} {
	var req dbrequest.DBRequest
	if err := json.Unmarshal(log.Data, &req); err != nil {
		s.logger.Error("raft.Apply(): failed to unmarshal request", zap.Error(err))
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
		s.logger.Warn("raft.Apply(): should not be applying reads", zap.String("key", req.Key))
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
	}
	return nil
}

// snapshotting is an optimization that we should actually implement later
type snapshotNoop struct{}

func (sn snapshotNoop) Persist(_ raft.SnapshotSink) error { return nil }
func (sn snapshotNoop) Release()                          {}
func (store *Store) Snapshot() (raft.FSMSnapshot, error) {
	return snapshotNoop{}, nil
}
