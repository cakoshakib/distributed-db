package boltstore

import (
	"encoding/binary"
	"errors"

	"github.com/hashicorp/raft"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var (
	LogsKey = []byte("logs")
	ConfKey = []byte("conf")
)

type BoltStore struct {
	db     *bolt.DB
	logger *zap.Logger
}

func NewBoltStore(db *bolt.DB, logger *zap.Logger) *BoltStore {
	return &BoltStore{
		db:     db,
		logger: logger,
	}
}

func (bls *BoltStore) GetLog(index uint64, log *raft.Log) error {
	return bls.db.View(func(tx *bolt.Tx) error {
		logsBucket := tx.Bucket(LogsKey)
		logBytes := logsBucket.Get(itob(index))
		if logBytes == nil {
			return raft.ErrLogNotFound
		}
		gotLog, err := LogFromPB(logBytes)
		if err != nil {
			return err
		}
		log.Index = gotLog.Index
		log.Term = gotLog.Term
		log.Type = gotLog.Type
		log.Data = gotLog.Data
		log.Extensions = gotLog.Extensions
		log.AppendedAt = gotLog.AppendedAt

		return nil
	})
}

func (bls *BoltStore) StoreLog(log *raft.Log) error {
	return bls.StoreLogs([]*raft.Log{log})
}

func (bls *BoltStore) StoreLogs(logs []*raft.Log) error {
	return bls.db.Batch(func(tx *bolt.Tx) error {
		logsBucket := tx.Bucket(LogsKey)
		for _, log := range logs {
			data, err := LogToPB(log)
			if err != nil {
				return err
			}
			if err := logsBucket.Put(itob(log.Index), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bls *BoltStore) DeleteRange(min, max uint64) error {
	return bls.db.Batch(func(tx *bolt.Tx) error {
		logsBucket := tx.Bucket(LogsKey)
		for i := min; i <= max; i++ {
			if err := logsBucket.Delete(itob(i)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bls *BoltStore) FirstIndex() (uint64, error) {
	var index uint64
	err := bls.db.View(func(tx *bolt.Tx) error {
		logsBucket := tx.Bucket(LogsKey)
		cursor := logsBucket.Cursor()
		k, _ := cursor.First()
		if k == nil {
			index = 0
			return nil
		}
		index = btoi(k)
		return nil
	})
	return index, err
}

func (bls *BoltStore) LastIndex() (uint64, error) {
	var index uint64
	err := bls.db.View(func(tx *bolt.Tx) error {
		logsBucket := tx.Bucket(LogsKey)
		cursor := logsBucket.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			index = 0
			return nil
		}
		index = btoi(k)
		return nil
	})
	return index, err
}

// Util functions
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// functions implemented for stable store below:

func (bls *BoltStore) Set(k, v []byte) error {
	tx, err := bls.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(ConfKey)
	if err := bucket.Put(k, v); err != nil {
		return err
	}

	return tx.Commit()
}

func (bls *BoltStore) Get(k []byte) ([]byte, error) {
	tx, err := bls.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(ConfKey)
	val := bucket.Get(k)

	if val == nil {
		return nil, errors.New("not found")
	}
	return append([]byte(nil), val...), nil
}

func (bls *BoltStore) SetUint64(key []byte, val uint64) error {
	return bls.Set(key, itob(val))
}

func (bls *BoltStore) GetUint64(key []byte) (uint64, error) {
	val, err := bls.Get(key)
	if err != nil {
		return 0, err
	}
	return btoi(val), nil
}
