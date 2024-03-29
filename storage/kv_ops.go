package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func add_kv_to_file(user string, table string, key string, value interface{}) error {
	switch value.(type) {
	case string:
		// value is string

	case map[string]interface{}:
		// value is JSON
	case []interface{}:
		// value is list
	case int:
		// value is a number
	}

	// get UserTable of given request, or create it
	// - get User struct
	// add key value pair to the table
	// persist data in file
	// return
	return nil
}

func read_kv_from_file(user string, table string, key string) (interface{}, error) {
	table_path := table_path(user, table)
	file, err := ioutil.ReadFile(table_path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s", table_path)
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling data %s", table_path)
	}

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
