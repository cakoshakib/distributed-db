package storage

import (
	"fmt"
	"os"
)

func add_user_to_file(user string) {
	file_path := "data/" + user
	if _, err := os.Stat(file_path); err == nil {
		// user already exists
	} else {
		err := os.Mkdir(file_path, 0755)
		if err != nil {
			fmt.Println("error creating file " + file_path)
		}
	}
}

func add_table_to_file(user string, table string) {
}

func AddUser(user string) {
	fmt.Println("Attempting to add user " + user)
	add_user_to_file(user)
}
