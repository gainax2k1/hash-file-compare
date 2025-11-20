package hashfile

/*
	goal: to recieve a file/filename? (prolly filename) and return the
	sha256 value from it.

	maybe instead use CLI command "sha256sum <filename>" instead?

*/
import (
	"crypto/sha256"
	"fmt"
	"hash"
)

type hasher struct {
	hash hash.Hash
}

func newHasher() *hasher { //returns pointer to new hasher
	var newHash hasher
	newHash.hash = sha256.New()
	return &newHash
}

func (h *hasher) Write(msg string) (int, error) {
	/*
	   type Writer interface {
	   	Write(p []byte) (n int, err error)
	   }

	*/
	// writes data to the hasher. accepts a string and casts to []byte
	bMsg := []byte(msg)
	n, err := h.hash.Write(bMsg)
	if err != nil {
		return n, err // i believe this is what the lesson requested?
	}
	return n, nil
}

func (h *hasher) GetHex() string {
	//method on a pointer to hasher. get's the hash value of the data written to the hasher
	//encodes  the hash value a lowercase hex string and returns it
	hashValue := h.hash.Sum(nil)
	return fmt.Sprintf("%x", hashValue)

}

/*
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

	filename := os.Args[1]

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()

	// Create a new SHA256 hash object
	hasher := sha256.New()

	// Copy file data into the hasher
	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatalf("Error hashing file: %v\n", err)
	}

	// Compute the SHA256 checksum
	checksum := hasher.Sum(nil)

	// Convert the checksum to hexadecimal format
	checksumHex := hex.EncodeToString(checksum)

	// Print the SHA256 checksum
	fmt.Printf("%s  %s\n", checksumHex, filename)
}



    The program checks if a filename is provided as a command-line argument.
    It opens the specified file for reading.
    A SHA256 hash object is created.
    The contents of the file are copied to the hasher.
    The SHA256 checksum is computed and converted to a hexadecimal string.
    Finally, it prints the checksum along with the filename.

*/
