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
func user_path(dataDir string, user string) string {
	return fmt.Sprintf("%s%s", dataDir, user)
}

func table_path(dataDir string, user string, table string) string {
	return fmt.Sprintf("%s/%s.json", user_path(dataDir, user), table)
}
