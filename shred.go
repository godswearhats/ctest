package shred

// Implement a Shred(path) function that will overwrite the given file (e.g. “randomfile”) 3 times with random data and delete the file afterwards. Note that the file may contain any type of data.

import (
	"errors"
	"io/fs"
	"os"
)

func Shred(path string) (bool, error) {
	if !fs.ValidPath(path) {
		return false, errors.New("Path is invalid")
	}
	_, err := os.Open(path)
	if err != nil {
		return false, err
	}
	return true, nil
}
