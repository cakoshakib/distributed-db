package storage

import (
	"fmt"
	"os"
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

func AddTable(user string, table string) {
	fmt.Printf("attempting to add table %s to user %s\n", table, user)
	if err := add_table_to_file(user, table); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("created table %s to user %s\n", table, user)
}

func DeleteTable(user string, table string) {
	fmt.Printf("attempting to remove table %s to user %s\n", table, user)
	if err := remove_table_from_file(user, table); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("removed table %s to user %s\n", table, user)

}
