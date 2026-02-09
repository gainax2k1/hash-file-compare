package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

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
		returnedMap, err := walkDir.WalkDir(os.Args[2])
		if err != nil {
			log.Fatalf("Error walking directory: %v\n", err)
		}

		// Print hash files for debugging purposes
		/*
			for hash := range returnedMap {
				for _, path := range returnedMap[hash] {
					fmt.Printf("Hash: %s\nFiles: %v\n", hash, path)
				}
			}*/

		// Display duplicate files for debugging purposes
		fmt.Println("Printing duplicate files:")
		displayDupicateFiles(returnedMap)
		return
	}
	if os.Args[1] == "-REMOVE" {
		// verify there's a directory path argument
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s -REMOVE <directory_path>\n", os.Args[0])
		}

		// call WalkDir with the provided directory path
		returnedMap, err := walkDir.WalkDir(os.Args[2])
		if err != nil {
			log.Fatalf("Error walking directory: %v\n", err)
		}

		// Remove duplicate files
		deleteDuplicateFiles(returnedMap)
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
			fmt.Printf("Hash: %s", hash)
			fmt.Println("Files:")
			for _, path := range paths {
				fmt.Printf(" - %s\n", path)

			}
		}
	}
}

func getTrashPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get current user: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		return "C:\\$Recycle.Bin\\", nil
	case "darwin":
		return filepath.Join("/Users", usr.Username, ".Trash/"), nil
	case "linux":
		return filepath.Join("/home", usr.Username, ".local/share/Trash/files/"), nil
	default:
		return "", fmt.Errorf("unsupported OS")
	}
}

func trashDuplicateFiles(hashMap map[string][]string) error {
	//Get username for trash path
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to get current user: %v", err)
	}

	// Define the trash path based on the OS
	var trashPath string

	switch runtime.GOOS {
	case "windows":
		trashPath = "C:\\$Recycle.Bin\\"
	case "darwin":
		trashPath = "/Users" + usr.Username + ".Trash/"
	case "linux":
		trashPath = "/home" + usr.Username + ".local/share/Trash/files/"
	default:
		return fmt.Errorf("unsupported OS for trashing files")
	}

	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and delete the rest
			for i := 1; i < len(paths); i++ {
				// Move the file to the trash, adding trashPath to the file name
				err := os.Rename(paths[i], trashPath+"/"+paths[i])

				if err != nil {
					log.Printf("Error moving to trash file %s: %v\n", paths[i], err)
				} else {
					fmt.Printf("Trashed file: %s\n", paths[i])
				}
			}
		}
	}
}

func deleteDuplicateFiles(hashMap map[string][]string) {
	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and delete the rest
			for i := 1; i < len(paths); i++ {
				err := os.Remove(paths[i])
				if err != nil {
					log.Printf("Error deleting file %s: %v\n", paths[i], err)
				} else {
					fmt.Printf("Deleted duplicate file: %s\n", paths[i])
				}
			}
		}
	}
}
