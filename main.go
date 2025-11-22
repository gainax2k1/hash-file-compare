package main

import (
	"fmt"
	"log"
	"os"

	hashfile "github.com/gainax2k1/hash-file-compare/hashFile"
)

func main() {
	fmt.Println("Find duplicate files by hash value")
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

	filename := os.Args[1]

	fileHashValue, err := hashfile.HashFile(filename)
	if err != nil {
		//handle error here!!! *********************
	}
	fmt.Println(fileHashValue)

}
