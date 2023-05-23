package shred

import (
	"os"
	"testing"
)

func TestShouldReturnErrorIfPathInvalid(t *testing.T) {
	result, err := Shred("this/../path/../is/../invalid.txt")
	if err == nil || result {
		t.Errorf("No error returned for invalid path")
	}
}

func TestShouldReturnErrorIfFileDoesNotExist(t *testing.T) {
	result, err := Shred("non-existent-file.txt")
	if err == nil || result {
		t.Errorf("No error returned for non-existent file")
	}
}

func TestShouldReturnErrorIfFileIsDirectory(t *testing.T) {
	path := "test_dir_delete_me"
	os.Mkdir(path, 0755)
	result, err := Shred(path)
	if err == nil || result {
		t.Errorf("No error returned when asked to shred directory")
	}
	os.Remove(path)
}
