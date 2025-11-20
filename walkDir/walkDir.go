package walkdir

// goal: walk through all the files in a directory
// then call the hashFile mod with each file
// then each returned hash should be stored in a map
// and then returned?
// maybe break this up differently.
// currently, this is doing both the walk and the sha256 hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <directory_path>\n", os.Args[0])
	}

	dir := os.Args[1]
	hashMap := make(map[string][]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Process only files, ignore directories
		if !info.IsDir() {
			hashValue, err := hashFile(path)
			if err != nil {
				log.Printf("Error hashing file %s: %v\n", path, err)
				return nil // Continue with next file
			}
			hashMap[hashValue] = append(hashMap[hashValue], path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v\n", err)
	}

	// Print duplicate files
	for hash, files := range hashMap {
		if len(files) > 1 {
			fmt.Printf("Hash: %s\nFiles: %v\n", hash, files)
		}
	}
}

// hashFile computes the SHA256 hash for a given file
func hashFile(filename string) (string, error) {
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
