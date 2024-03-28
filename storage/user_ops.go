package storage

import (
	"fmt"
	"os"
)

func add_user_to_file(user string) error {
	user_path := fmt.Sprintf("data/%s", user)
	if !file_exists(user_path) {
		err := os.Mkdir(user_path, 0755)
		if err != nil {
			return fmt.Errorf("error creating file %s", user_path)
		}
	}
	return nil
}

func add_table_to_file(user string, table string) error {
	user_path := fmt.Sprintf("data/%s", user)
	table_path := fmt.Sprintf("%s/%s.json", user_path, table)
	// TODO: improve by making list of error msgs somewhere
	if !file_exists(user_path) {
		// User does not exist
		return fmt.Errorf("user does not exist %s", user)
	}
	if !file_exists(table_path) {
		// create table
		if _, err := os.Create(table_path); err != nil {
			return fmt.Errorf("error creating file %s", table_path)
		}
	}
	return nil
}

func AddUser(user string) {
	fmt.Println("attempting to add user " + user)
	if err := add_user_to_file(user); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("created user " + user)
}

func AddTable(user string, table string) {
	fmt.Printf("attempting to add table %s to user %s\n", table, user)
	if err := add_table_to_file(user, table); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("created table %s to user %s\n", table, user)
}
