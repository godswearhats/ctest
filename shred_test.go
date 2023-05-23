package shred

import (
	"os"
	"testing"
)

func TestShredShouldReturnErrorIfPathInvalid(t *testing.T) {
	result, err := Shred("this/../path/../is/../invalid.txt")
	if err == nil || result {
		t.Errorf("No error returned for invalid path")
	}
}

func TestShredShouldReturnErrorIfFileDoesNotExist(t *testing.T) {
	result, err := Shred("non-existent-file.txt")
	if err == nil || result {
		t.Errorf("No error returned for non-existent file")
	}
}

func TestShredReturnErrorIfFileIsDirectory(t *testing.T) {
	path := "test_dir_delete_me"
	os.Mkdir(path, 0755)
	result, err := Shred(path)
	if err == nil || result {
		t.Errorf("No error returned when asked to shred directory")
	}
	os.Remove(path)
}

func TestOverwriteShouldReplaceContentsOfTextFile(t *testing.T) {
	path := "test_file_delete_me.txt"
	expected := []byte("lorem ipsum blah blah blah")
	err := os.WriteFile(path, expected, 0644)
	if err != nil {
		panic(err)
	}

	err = Overwrite(path)
	if err != nil {
		t.Errorf("Overwrite returned error when replacing text file")
	}

	actual, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if SlicesEqual(expected, actual) {
		t.Errorf("Overwrite did not change file contents")
	}
	os.Remove(path)
}

// can't compare slice equality with regular '==' so rolled my own comparator
func SlicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
