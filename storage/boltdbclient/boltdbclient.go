package boltdbclient

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func NewBoltDB(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("bolt db could not be opened with err: %s", err)
	}

	if err := initBoltDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initBoltDB(db *bolt.DB) error {
	// Bucket creation
	return db.Update(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte("logs"))
		if err != nil {
			return fmt.Errorf("initBoltDB(): create bucket err: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("conf"))
		if err != nil {
			return fmt.Errorf("initBoltDB(): create bucket err: %s", err)
		}
		return nil
	})
}
