package filesys

import (
	"os"
)

func PathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func IsDirectory(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func IsFile(filename string) bool {
	return !IsDirectory(filename)
}
