package walkdir

// goal: walk through all the files in a directory
// then call the hashFile mod with each file
// then each returned hash should be stored in a map
// and then returned?
// maybe break this up differently.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	hashFile "github.com/gainax2k1/hash-file-compare/hashFile"
)

func WalkDir(dir string) {
	// map to store hash values and corresponding file paths
	hashMap := make(map[string][]string)

	// Walk through the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//need to correct this, not ignore dirs
		// Process only files, ignore directories
		if !info.IsDir() {
			hashValue, err := hashFile.HashFile(path)
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

	// Print hash files
	for hash, files := range hashMap {
		if len(files) > 1 {
			fmt.Printf("Hash: %s\nFiles: %v\n", hash, files)
		}
	}
}
