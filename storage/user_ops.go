package storage

import (
	"context"
	"fmt"
	"os"

	log "github.com/cakoshakib/distributed-db/commons"
	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"go.uber.org/zap"
)

// TODO: initialize data/ folder on server startup if it doesnt exist,
// otherwise, we have to manually create data/ before running funcs below
func add_user_to_file(user string) error {
	user_path := user_path(user)
	if !file_exists(user_path) {
		err := os.Mkdir(user_path, 0755)
		if err != nil {
			return fmt.Errorf("error creating dir %s", user_path)
		}
	}
	return nil
}

func remove_user_from_file(user string) error {
	user_path := user_path(user)
	if file_exists(user_path) {
		if err := os.RemoveAll(user_path); err != nil {
			return fmt.Errorf("error removing dir %s, may have tables", user_path)
		}
	}
	return nil
}

func DeleteUser(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	if err := remove_user_from_file(req.User); err != nil {
		logger.Error("storage: delete user error", zap.Error(err))
		return err
	}
	logger.Info("storage: delete user success", zap.String("user", req.User))
	return nil
}

func AddUser(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	if err := add_user_to_file(req.User); err != nil {
		logger.Error("storage: add user error", zap.Error(err))
		return err
	}
	logger.Info("storage: add user success", zap.String("user", req.User))
	return nil
}
