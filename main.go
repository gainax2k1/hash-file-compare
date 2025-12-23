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

	// check for -d flag here to call WalkDir with folder name ( i believe? )
	if os.Args[1] == "-d" {
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s -d <directory_path>\n", os.Args[0])
		}
		walkDir.WalkDir()
		return
	}
	// case switch, perhaps?

	// handles single file hash value check
	filename := os.Args[1]

	fileHashValue, err := hashfile.HashFile(filename)
	if err != nil {
		log.Fatalf("Error hashing file: %v\n", err)
	}

	fmt.Println(fileHashValue)

}
