package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	"go.uber.org/zap"
)

func read_json_data(user string, table string) map[string]interface{} {
	file, _ := ioutil.ReadFile(table_path(user, table))
	data := map[string]interface{}{}
	json.Unmarshal(file, &data)
	return data
}

func write_json_data(user string, table string, data []byte) error {
	return os.WriteFile(table_path(user, table), data, 0777)
}

func add_kv_to_file(user string, table string, key string, value interface{}) error {
	data := read_json_data(user, table)
	data[key] = value
	m, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data %v", data)
	}
	if err := write_json_data(user, table, m); err != nil {
		return fmt.Errorf("error writing json to file %v", data)
	}

	return nil
}

func remove_kv_from_file(user string, table string, key string) error {
	data := read_json_data(user, table)
	delete(data, key)
	m, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data %v", data)
	}
	if err := write_json_data(user, table, m); err != nil {
		return fmt.Errorf("error writing json to file %v", data)
	}
	return nil
}

func read_kv_from_file(user string, table string, key string) (interface{}, error) {
	data := read_json_data(user, table)
	value, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return value, nil
}

func (s *Store) ReadKV(req dbrequest.DBRequest) (string, error) {
	data, err := read_kv_from_file(req.User, req.Table, req.Key)
	if err != nil {
		s.logger.Error("storage: readkv error", zap.Error(err))
		return "", err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("storage: error marshaling data to json", zap.Error(err))
		return "", err
	}

	s.logger.Info("storage: readkv success", zap.String("key", req.Key), zap.ByteString("value", jsonData))
	return string(jsonData), nil
}

func (s *Store) AddKV(req dbrequest.DBRequest) error {

	var value interface{}
	err := json.Unmarshal([]byte(req.Value), &value)
	if err != nil {
		// assume value is not JSON
		value = req.Value
	}

	if err := add_kv_to_file(req.User, req.Table, req.Key, value); err != nil {
		s.logger.Error("storage: addkv error", zap.Error(err))
		return err
	}
	// same as above
	s.logger.Info("storage: addkv success", zap.String("key", req.Key))
	return nil
}

func (s *Store) RemoveKV(req dbrequest.DBRequest) error {
	err := remove_kv_from_file(req.User, req.Table, req.Key)
	if err != nil {
		s.logger.Error("storage: removekv error", zap.Error(err))
		fmt.Printf("err: %v\n", err)
		return err
	}
	s.logger.Info("storage: removekv success", zap.String("key", req.Key))
	return nil
}
