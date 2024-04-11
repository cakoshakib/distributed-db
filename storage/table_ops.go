package storage

import (
	"fmt"
	"os"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"go.uber.org/zap"
)

func add_table_to_file(user string, table string) error {
	user_path := user_path(user)
	table_path := table_path(user, table)
	// TODO: improve by making error structs
	if !file_exists(user_path) {
		return fmt.Errorf("user does not exist %s", user)
	}
	if !file_exists(table_path) {
		if _, err := os.Create(table_path); err != nil {
			return fmt.Errorf("error creating file %s", table_path)
		}
	}
	return nil
}

func remove_table_from_file(user string, table string) error {
	table_path := table_path(user, table)
	// TODO: improve by making error structs
	if file_exists(table_path) {
		if err := os.Remove(table_path); err != nil {
			return fmt.Errorf("error removing file %s", table_path)
		}
	}
	return nil
}

func (s *Store) AddTable(req dbrequest.DBRequest) error {
	if err := add_table_to_file(req.User, req.Table); err != nil {
		s.logger.Error("storage: add table error", zap.Error(err))
		return err
	}
	s.logger.Info("storage: add table success", zap.String("user", req.User), zap.String("table", req.Table))
	return nil
}

func (s *Store) DeleteTable(req dbrequest.DBRequest) error {
	if err := remove_table_from_file(req.User, req.Table); err != nil {
		s.logger.Info("storage: delete table error", zap.Error(err))
		return err
	}
	s.logger.Info("storage: delete table success", zap.String("user", req.User), zap.String("table", req.Table))
	return nil
}
