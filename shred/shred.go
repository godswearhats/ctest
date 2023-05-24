package shred

// Implement a Shred(path) function that will overwrite the given file (e.g. “randomfile”) 3 times with random data and delete the file afterwards. Note that the file may contain any type of data.

// In a few lines briefly discuss the possible use cases for such a helper function as well as advantages and drawbacks of addressing them with this approach.

// It seems like the use case here is to try to make it so deletion of a file from the filesystem is a bit more secure. Often when files are removed they are comparatively easy to restore, so overwriting them multiple times with random data makes it a bit more difficult to restore the data that was in the file originally. I imagine this will depend a lot on the storage medium that the file is actually stored on, and how the filesystem stores the data.

import (
	"errors"
	"io/fs"
	"math"
	"math/rand"
	"os"
)

const chunkSize int = 1 << 20 // using a chunk size of 1Mb as a reasonable default, this could be altered depending on use case

func Shred(path string) error {
	size, err := ValidatePathAndFindSize(path)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		err = Overwrite(path, size)
		if err != nil {
			return err
		}
	}

	return os.Remove(path)
}

func ValidatePathAndFindSize(path string) (int, error) {
	// Validate the path leads to an existing regular file
	if !fs.ValidPath(path) {
		return 0, errors.New("Specified path is invalid")
	}
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if fileInfo.IsDir() {
		return 0, errors.New("Specified file is a directory")
	}

	// grab the file size and return it
	size := int(fileInfo.Size())
	return size, nil
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
