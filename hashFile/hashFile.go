package hashfile

/*
	goal: to recieve a file/filename? (prolly filename) and return the
	sha256 value from it.

	maybe instead use CLI command "sha256sum <filename>" instead?

*/
import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// HashFile computes the SHA256 hash for a given file, returning the hash as a hex string.

func HashFile(filename string) (string, error) {
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
