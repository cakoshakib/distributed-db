package storage

import (
	"fmt"
	"os"
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

func DeleteUser(user string) {
	fmt.Println("attempting to remove user " + user)
	if err := remove_user_from_file(user); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println("removed user " + user)
}

func AddUser(user string) {
	fmt.Println("attempting to add user " + user)
	if err := add_user_to_file(user); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println("created user " + user)
}
