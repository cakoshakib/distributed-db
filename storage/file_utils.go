package storage

import "os"

func file_exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
