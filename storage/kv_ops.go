package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cakoshakib/distributed-db/commons/dbrequest"
	log "github.com/cakoshakib/distributed-db/commons"
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
	return data[key], nil
}

// TODO add returning KV w/ interface lol
func ReadKV(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	data, err := read_kv_from_file(req.user, req.table, req.key)
	if err != nil {
		logger.Error("storage: readkv error", zap.Error(err))
		return err
	}
	// add data to this log after implementing data -> string operation
	logger.Info("storage: readkv success", zap.String("key", req.key)) 
	return nil
}

func AddKV(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	if err := add_kv_to_file(req.user, req.table, req.key, req.value); err != nil {
		zap.Error("storage: addkv error", zap.Error(err))
		return err
	}
	// same as above
	zap.Info("storage: addkv success", zap.String("key", req.key))
	return nil
}

func RemoveKV(ctx context.Context, req dbrequest.DBRequest) error {
	logger := log.LoggerFromContext(ctx)
	err := remove_kv_from_file(req.user, req.table, req.key)
	if err != nil {
		zap.Error("storage: removekv error", zap.Error(err))
		fmt.Printf("err: %v\n", err)
		return err
	}
	zap.Info("storage: removekv success", zap.String("key", req.key))
	return nil
}
