package main

import (
	"archive/tar"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Open the archive
	var path string
	argAmount := len(os.Args)
	if argAmount < 2 {
		log.Fatal(`There were not enough arguments.
	 This command requires a path to be given.`)
	} else if argAmount > 2 {
		// TODO: Allow multiple paths
		log.Fatal(`There are too many arguments.
Only a single path at a time is allowed`)
	}

	args := os.Args[1]
	fmt.Printf("here are some args:\n%s", args)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory")
	}
	// Need to look for a file that includes a file with a mets.
	// The mets will have a uuid: mets.<uuid>.xml
	dataPath := filepath.Join(cwd, path, "data")
	cwd = dataPath
	dir, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatalf("This path does not exist.")
	}
	var entry os.DirEntry
	for _, dEntry := range dir {
		if strings.Contains(dEntry.Name(), "mets") && !dEntry.IsDir() {
			entry = dEntry
			break
		}
	}
	mets, err := os.Open(entry.Name())
	defer func {
		err := mets.Close()
		if err != nil {
			log.Fatal("Failed to close the mets file.")
		}
	}()
	if err != nil {
		log.Fatalf("Could not open %s", entry.Name())
	}
	// We want to read the file
	r := tar.NewReader(mets)
	_ = r

	// TODO: We need to loop over the archive to find the mets file.
	// Which should be in the top level of the data directory.

}
