package hashfile

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

const PARTIALBYTELENGTH = 4096 //bytes

// computeHash handles the actual hashing logic for both hash sizes
func computeHash(r io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func FullHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return computeHash(file)
}

func PartialHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// io.LimitReader ensures we only read up to the limit
	// but computeHash will still work if the file is smaller
	limitedReader := io.LimitReader(file, PARTIALBYTELENGTH)
	return computeHash(limitedReader)
}
