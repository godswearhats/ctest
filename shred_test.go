package shred

import "testing"

func TestShouldReturnErrorIfPathInvalid(t *testing.T) {
	_, err := Shred("nonexistent_file.txt")
	if err == nil {
		t.Errorf("No error returned for non-existent file")
	}
}
