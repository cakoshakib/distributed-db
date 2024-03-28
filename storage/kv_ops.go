package storage

import (
	"fmt"
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

func test() {
	fmt.Println("hello world")
}
