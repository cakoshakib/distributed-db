package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func ReadKV(user string, table string, key string) {
	fmt.Println("attempting to read key", key)
	data, err := read_kv_from_file(user, table, key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("got data: %v\n", data)
}

func AddKV(user string, table string, key string, value interface{}) {
	fmt.Println("attempting to add key, value", key, value)
	if err := add_kv_to_file(user, table, key, value); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("wrote kv pair", key, value)
}

func RemoveKV(user string, table string, key string) {
	fmt.Println("attempting to read key", key)
	err := remove_kv_from_file(user, table, key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("removed key %s\n", key)

}
