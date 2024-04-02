package storage

import (
	"fmt"
	"os"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	log "github.com/cakoshakib/distributed-db/commons"
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

func AddTable(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	if err := add_table_to_file(req.user, req.table); err != nil {
		logger.Error("storage: add table error", zap.Error(err))
		return err
	}
	logger.Info("storage: add table success", zap.String("user", req.user), zap.String("table", req.table))
	return nil
}

func DeleteTable(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	if err := remove_table_from_file(user, table); err != nil {
		logger.Info("storage: delete table error", zap.Error(err))
		return err
	}
	logger.Info("storage: delete table success", zap.String("user", req.user), zap.String("table", req.table))
	return nil
}
