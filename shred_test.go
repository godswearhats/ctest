package shred

import "testing"

func TestShouldReturnErrorIfPathInvalid(t *testing.T) {
	result, err := Shred("this/../path/../is/../invalid.txt")
	if err == nil || result {
		t.Errorf("No error returned for invalid path")
	}
}
