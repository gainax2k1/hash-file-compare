package walkdir

// goal: walk through all the files in a directory
// then call the hashFile mod with each file
// then each returned hash should be stored in a map
// and then returned?
// maybe break this up differently.

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	hashFile "github.com/gainax2k1/hash-file-compare/hashFile"
)

//takes in a directory path, returns a map of hash values (of each file), and slice of the file paths that correspond to each hash value

func WalkDir(dir string) (map[string][]string, error) {
	// map to store hash values and corresponding file paths
	hashMap := make(map[string][]string)

	//refactoring using WalkDir function (replaced deprecated Walk function)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Process only files, ignore directories
		if !d.IsDir() {
			hashValue, err := hashFile.HashFile(path)
			// fmt.Println(hashValue, path)
			if err != nil {
				log.Printf("Error hashing file %s: %v\n", path, err)
				return nil // Continue with next file
			}

			hashMap[hashValue] = append(hashMap[hashValue], path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error walking: %w", err)
	}
	// Check if any files were processed
	if len(hashMap) == 0 {
		return nil, errors.New("no files found in the directory")
	}

	// Return the map of hash values and file paths
	return hashMap, nil
}
