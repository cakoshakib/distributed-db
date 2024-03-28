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
			return fmt.Errorf("error creating dir %s", user_path)
		}
	}
	return nil
}

func remove_user_from_file(user string) error {
	user_path := fmt.Sprintf("data/%s", user)
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
	}
	fmt.Println("removed user " + user)
}

func AddUser(user string) {
	fmt.Println("attempting to add user " + user)
	if err := add_user_to_file(user); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("created user " + user)
}
