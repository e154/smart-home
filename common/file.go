package common

import (
	"path/filepath"
	"unicode/utf8"
	"os"
)

func GetFileSize(name string) (int64, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func GetFullPath(name string) string {

	const (
		dataDir  = "./data"
		fileStoragePath = "./file_storage"
		depth = 3
	)

	dir := filepath.Join(dataDir, fileStoragePath)

	for i := 0; i < depth; i++ {
		dir = filepath.Join(dir, name[i*3:(i+1)*3])
	}

	return dir
}

func GetLinkPath(name string) string {

	dir := "/upload"

	count := utf8.RuneCountInString(name)
	if count < 9 {
		return filepath.Join(dir, name)
	}

	const depth = 3

	for i := 0; i < depth; i++ {
		dir = filepath.Join(dir, name[i*3:(i+1)*3])
	}

	return filepath.Join(dir, name)
}