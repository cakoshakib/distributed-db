package storage

import (
	"fmt"
	"os"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"go.uber.org/zap"
)

// TODO: initialize data/ folder on server startup if it doesnt exist,
// otherwise, we have to manually create data/ before running funcs below
func add_user_to_file(user_path string) error {
	if !file_exists(user_path) {
		err := os.Mkdir(user_path, 0755)
		if err != nil {
			return fmt.Errorf("error creating dir %s", user_path)
		}
	}
	return nil
}

func remove_user_from_file(user_path string) error {
	if file_exists(user_path) {
		if err := os.RemoveAll(user_path); err != nil {
			return fmt.Errorf("error removing dir %s, may have tables", user_path)
		}
	}
	return nil
}

func (s *Store) DeleteUser(req dbrequest.DBRequest) error {
	if err := remove_user_from_file(user_path(s.dataDir, req.User)); err != nil {
		s.logger.Error("storage: delete user error", zap.Error(err))
		return err
	}
	s.logger.Info("storage: delete user success", zap.String("user", req.User))
	return nil
}

func (s *Store) AddUser(req dbrequest.DBRequest) error {
	if err := add_user_to_file(user_path(s.dataDir, req.User)); err != nil {
		s.logger.Error("storage: add user error", zap.Error(err))
		return err
	}
	s.logger.Info("storage: add user success", zap.String("user", req.User))
	return nil
}
