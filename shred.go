package shred

// Implement a Shred(path) function that will overwrite the given file (e.g. “randomfile”) 3 times with random data and delete the file afterwards. Note that the file may contain any type of data.

import (
	"errors"
	"io/fs"
	"os"
)

func Shred(path string) (bool, error) {
	// Validate the path leads to an existing regular file
	if !fs.ValidPath(path) {
		return false, errors.New("Specified path is invalid")
	}
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if fileInfo.IsDir() {
		return false, errors.New("Specified file is a directory")
	}

	return true, nil
}
