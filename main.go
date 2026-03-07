package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gainax2k1/hashcomparefiles/internal/logger"
	walkdir "github.com/gainax2k1/hashcomparefiles/internal/walkdir"
)

type Config struct {
	Path             string
	Trash            bool
	Delete           bool
	Verbose          bool
	LogPath          string
	ShowPreHashCount bool
}

func main() {
	// Define flags and parse
	trashFlag := flag.Bool("trash", false, "Trash duplicate files instead of just listing")
	deletFlag := flag.Bool("delete", false, "Delete duplicate files instead of just listing")
	logFlag := flag.String("log", "none", "Log path, or 'default' for current directory")
	showPreHashCountFlag := flag.Bool("p", false, "Show Pre-hash file count (Potentially usefull for large runs, but now hits storage twice)")
	verboseFlag := flag.Bool("v", false, "verbose mode,")

	flag.Parse()

	// Identify all paths to process (pipe or args)
	var targets []string

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if path := strings.TrimSpace(scanner.Text()); path != "" {
				targets = append(targets, path)
			}
		}
	} else {
		// Use command line arguments if no pipe
		targets = flag.Args()
	}

	// Validate targets
	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <path>\n", os.Args[0])
		os.Exit(1)
	}

	// Create config struct with parsed values
	config := Config{
		Trash:            *trashFlag,
		Delete:           *deletFlag,
		Verbose:          *verboseFlag,
		LogPath:          *logFlag,
		ShowPreHashCount: *showPreHashCountFlag,
	}

	// All output will be done through the logger, writing to file and/or screen based on config
	logger, err := logger.NewLogger(config.LogPath, config.Verbose)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer logger.Close()

	err = process(targets, config, logger)
	if err != nil {
		logger.Error("Error running directory mode: %v", err)
	}

}

func process(targets []string, config Config, logger *logger.Logger) error {

	masterMap := make(map[string][]walkdir.FileInfo)
	totalCount := 0
	runHash := false

	if config.ShowPreHashCount {
		for _, path := range targets {
			// Run the walk for each path and merge results into masterMap
			_, count, err := walkdir.WalkDir(path, logger, config.Verbose, runHash)
			if err != nil {
				logger.Error("Skipping %s due to error: %v", path, err)
				continue // Keep going with other targets!
			}

			totalCount += count
		}
		logger.Log("Total files to process: %d", totalCount)
		totalCount = 0 // reset for hashing run
	}

	runHash = true
	for _, path := range targets {
		// Run the walk for each path and merge results into masterMap
		dirMap, count, err := walkdir.WalkDir(path, logger, config.Verbose, runHash)
		if err != nil {
			logger.Error("Skipping %s due to error: %v", path, err)
			continue // Keep going with other targets!
		}

		totalCount += count
		// Merge dirMap into masterMap
		for hash, files := range dirMap {
			masterMap[hash] = append(masterMap[hash], files...)
		}
	}

	if config.Trash {
		if err := trashDuplicateFiles(masterMap, logger); err != nil {
			return fmt.Errorf("Error trashing duplicate files: %w", err)
		}
	} else if config.Delete {
		if err := deleteDuplicateFiles(masterMap, logger); err != nil {
			return fmt.Errorf("Error deleting duplicate files: %w", err)
		}
	} else {
		// just list duplicates, do nothing else

		displayHashMap(logger, masterMap, totalCount, config)
	}

	return nil
}

func displayHashMap(logger *logger.Logger, hashMap map[string][]walkdir.FileInfo, count int, config Config) {
	for hash, paths := range hashMap {
		if config.Verbose {
			// if verbose, print all files, even if not duplicates, and include file sizes
			logger.Log("Hash: %s", hash)

			for _, path := range paths {
				logger.Log(" - %s size: %d", path.FilePath, path.FileSize)
			}

		} else { // if not verbose, just print instances with duplicates

			if len(paths) > 1 {
				logger.Log("Duplicate files with hash: %s", hash)
				for _, path := range paths {
					logger.Log(" - %s size: %d", path.FilePath, path.FileSize)

				}
			}
		}

	}
	logger.Log("Total files processed: %d", count)
}

func trashDuplicateFiles(hashMap map[string][]walkdir.FileInfo, logger *logger.Logger) error {
	//Get username for trash path
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to get current user: %v", err)
	}

	// Define the trash path based on the OS
	var trashPath, trashInfoDir string

	if runtime.GOOS == "linux" {

		trashPath = filepath.Join(usr.HomeDir, ".local/share/Trash/files/")
		trashInfoDir = filepath.Join(usr.HomeDir, ".local/share/Trash/info/")
		// Ensure the trash info directory exists
		if _, err := os.Stat(trashInfoDir); os.IsNotExist(err) {
			err := os.MkdirAll(trashInfoDir, 0755)
			if err != nil {
				return fmt.Errorf("Error creating trash info directory: %v", err)
			}
		}
	} else {
		trashPath = "trash"
		trashInfoDir = "trash"
		// Ensure the trash directory exists
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			err := os.Mkdir(trashPath, 0755)
			if err != nil {
				return err
			}
		}
	}

	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and trash the rest
			// *Future improvement*: iterate through duplicates and ask user which one to keep,
			//  or if they want to keep all, trash all, etc. For now, just keep the first one and trash the rest.

			for i := 1; i < len(paths); i++ {

				// Create a unique name for the file in the trash to avoid conflicts
				ext := filepath.Ext(paths[i].FilePath)
				name := strings.TrimSuffix(filepath.Base(paths[i].FilePath), ext)
				enumeratedName := fmt.Sprintf("%s_%d%s", name, i, ext)

				destPath := filepath.Join(trashPath, enumeratedName)
				src := paths[i].FilePath

				// Move the file to the trash, adding trashPath to the file name
				// First try to rename (move) the file, which is more efficient.
				err := os.Rename(paths[i].FilePath, destPath)
				if err != nil {
					// Rename failed, try copy + delete method as a fallback
					err = copyFile(src, destPath)
					if err != nil {
						logger.Error("Error copying file to trash %s: %v", paths[i].FilePath, err)
						return err
					}
					err = os.Remove(src)
					if err != nil {
						logger.Error("Error deleting original file after copying to trash %s: %v", paths[i].FilePath, err)
						return err
					}

					logger.Log("Trashed file (copy+delete): %s", paths[i].FilePath)
				} else {
					logger.Log("Trashed file: %s", paths[i].FilePath)
				}

				// Create .trashinfo file (to FreeDesktop spec) if on Linux in appropriate directory, non-Linux will place .trashinfo files
				// in the same directory as the trashed files for simplicity

				infoPath := filepath.Join(trashInfoDir, enumeratedName+".trashinfo")
				originalPath := paths[i].FilePath
				infoContent := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", url.PathEscape(originalPath), time.Now().Format("2006-01-02T15:04:05"))

				err = os.WriteFile(infoPath, []byte(infoContent), 0644)
				if err != nil {
					logger.Error("Error creating trash info file for %s: %v", paths[i].FilePath, err)
					return err
				}

			}
		}
	}
	return nil
}

func deleteDuplicateFiles(hashMap map[string][]walkdir.FileInfo, logger *logger.Logger) error {
	// Iterate through the hash map and delete duplicate files, keeping the first instance
	// *Future improvement*: iterate through duplicates and ask user which one to keep.
	for _, paths := range hashMap {
		if len(paths) > 1 {
			// Keep the first file and delete the rest
			for i := 1; i < len(paths); i++ {
				err := os.Remove(paths[i].FilePath)
				if err != nil {
					logger.Error("Error deleting file %s: %v", paths[i].FilePath, err)
				} else {
					logger.Log("Deleted duplicate file: %s", paths[i].FilePath)
				}
			}

		}
	}
	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
