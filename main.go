package main

import (
	"fmt"
	"log"
	"os"

	hashfile "github.com/gainax2k1/hash-file-compare/hashFile"
	walkDir "github.com/gainax2k1/hash-file-compare/walkDir"
)

func main() {
	fmt.Println("Find duplicate files by hash value")
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

	// check for -d flag here to call WalkDir
	if os.Args[1] == "-d" {
		// verify there's a directory path argument
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s -d <directory_path>\n", os.Args[0])
		}
		// call WalkDir with the provided directory path
		// Need to rework this!
		returnedMap, err := walkDir.WalkDir(os.Args[2])
		if err != nil {
			log.Fatalf("Error walking directory: %v\n", err)
		}
		// Print hash files
		for hash := range returnedMap {
			for _, path := range returnedMap[hash] {
				fmt.Printf("Hash: %s\nFiles: %v\n", hash, path)
			}
		}
		fmt.Println("Printing duplicate files:")
		displayDupicateFiles(returnedMap)
		return
	}

	// handles single file hash value check
	filename := os.Args[1]

	fileHashValue, err := hashfile.HashFile(filename)
	if err != nil {
		log.Fatalf("Error hashing file: %v\n", err)
	}

	fmt.Println(fileHashValue)

}

func displayDupicateFiles(hashMap map[string][]string) {
	for hash, paths := range hashMap {
		if len(paths) > 1 {
			fmt.Printf("Hash: %s\nFiles: %v\n", hash, paths)
		}
	}
}
