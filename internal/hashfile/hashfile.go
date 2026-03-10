package hashfile

/*
	recieves a filename and return the sha256 value (hex) from it.
	plan: use the crypto/sha256 package to compute the hash value of the file.
	steps:
	1. open the file using os.Open
	2. create a new sha256 hasher using sha256.New()
	3. copy the file contents to the hasher using io.Copy
	4. compute the hash value and return it as a hex string using hex.EncodeToString
*/

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

const PARTIALBYTELENGTH = 4096 //bytes

// computeHash handles the actual hashing logic for both hash sizesS
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

/* old hashing functions, replaced with helper design, to more uniformaly handle the file hashes.
func FullHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}

func PartialHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, PARTIALBYTELENGTH)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	hash := sha256.Sum256(buffer[:n])
	return hex.EncodeToString(hash[:]), nil
}
*/
