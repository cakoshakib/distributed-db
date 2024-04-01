package storage

import (
	"fmt"
	"os"
)

func file_exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// TODO: remove dependency on relative paths. 
// perhaps set a configurable data dir path w/ default to */distributeddb/data
func user_path(user string) string {
	return fmt.Sprintf("data/%s", user)
}

func table_path(user string, table string) string {
	return fmt.Sprintf("%s/%s.json", user_path(user), table)
}
