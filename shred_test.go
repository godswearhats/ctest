package shred

import (
	"os"
	"testing"
)

// Some test cases that are missing:
// - validate that it works for binary files like images etc. (it should)
// - benchmarks: the chunk size will likely make a significant difference in how the code performs
// - it would be nice to be able to validate that Overwrite was called three times, could do this with mocking
//   but given that I've not written any Go code before today, I figured I'd limit the scope :-)

func TestShredShouldRemoveFile(t *testing.T) {
	expected := []byte("lorem ipsum blah blah blah")
	path := CreateTextFile(expected)

	err := Shred(path)
	if err != nil {
		t.Errorf("Shred returned error when replacing text file")
	}

	if _, err := os.Stat(path); err == nil {
		t.Errorf("Shred did not remove file")
	}
}

func TestShredShouldReturnErrorIfPathInvalid(t *testing.T) {
	err := Shred("this/../path/../is/../invalid.txt")
	if err == nil {
		t.Errorf("No error returned for invalid path")
	}
}

func TestShredShouldReturnErrorIfFileDoesNotExist(t *testing.T) {
	err := Shred("non-existent-file.txt")
	if err == nil {
		t.Errorf("No error returned for non-existent file")
	}
}

func TestShredShouldReturnErrorIfFileIsDirectory(t *testing.T) {
	path := "test_dir_delete_me"
	os.Mkdir(path, 0755)
	err := Shred(path)
	if err == nil {
		t.Errorf("No error returned when asked to shred directory")
	}
}

func TestOverwriteShouldReplaceContentsOfTextFile(t *testing.T) {
	expected := []byte("lorem ipsum blah blah blah")
	path := CreateTextFile(expected)

	err := Overwrite(path, len(expected))
	if err != nil {
		t.Errorf("Overwrite returned error when replacing text file")
	}

	AssertFileIsDifferent(path, expected, t)
	os.Remove(path)
}

func TestOverwriteShouldReplaceContentsOfTextFileWithDifferentDataEachTime(t *testing.T) {
	zeroth := []byte("lorem ipsum blah blah blah")
	path := CreateTextFile(zeroth)

	err := Overwrite(path, len(zeroth))
	if err != nil {
		t.Errorf("Overwrite returned error when replacing text file")
	}

	first := AssertFileIsDifferent(path, zeroth, t)

	err = Overwrite(path, len(zeroth))
	if err != nil {
		t.Errorf("Overwrite returned error when replacing text file")
	}

	second := AssertFileIsDifferent(path, first, t)

	err = Overwrite(path, len(zeroth))
	if err != nil {
		t.Errorf("Overwrite returned error when replacing text file")
	}

	third := AssertFileIsDifferent(path, second, t)
	if SlicesEqual(third, zeroth) {
		t.Errorf("Overwrite did not change file contents")
	}

	os.Remove(path)
}

// some helper functions to keep tests more readable

func AssertFileIsDifferent(path string, expected []byte, t *testing.T) []byte {
	actual, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if SlicesEqual(expected, actual) {
		t.Errorf("Overwrite did not change file contents (expected: %v actual: %v)", expected, actual)
	}
	return actual
}

func CreateTextFile(contents []byte) string {
	path := "test_file_delete_me.txt"
	err := os.WriteFile(path, contents, 0644)
	if err != nil {
		panic(err) // keep it simple for test code, if there's an issue we should just barf
	}
	return path
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
