package shred

// Implement a Shred(path) function that will overwrite the given file (e.g. “randomfile”) 3 times with random data and delete the file afterwards. Note that the file may contain any type of data.

import (
	"errors"
	"io/fs"
	"math"
	"math/rand"
	"os"
)

const chunkSize int = 1 << 20 // using a chunk size of 1Mb as a reasonable default, this could be altered depending on use case

func Shred(path string) (bool, error) {
	err := ValidatePath(path)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ValidatePath(path string) error {
	// Validate the path leads to an existing regular file
	if !fs.ValidPath(path) {
		return errors.New("Specified path is invalid")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if fileInfo.IsDir() {
		return errors.New("Specified file is a directory")
	}
	return nil
}

func Overwrite(path string, size int) error {
	// open the file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// break the operation into reasonable size chunks
	for written := 0; written < size; written += chunkSize {
		nextChunkSize := math.Min(float64(size-written), float64(chunkSize))
		randomBytes := GetRandomBytes(int(nextChunkSize))
		numWritten, err := file.Write(randomBytes)
		if err != nil {
			return err
		}
		written += numWritten
	}

	file.Close()

	return err
}

func GetRandomBytes(size int) []byte {
	var result []byte
	result = make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = byte(rand.Intn(256)) // bytes are unsigned ints, so we return a value between 0 and 255
	}
	return result
}
