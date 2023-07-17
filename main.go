package main

import (
	"archive/tar"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	mets_path := filepath.Join(cwd, path, "data", "mets.xml")
	_ = mets_path
	file, err := os.Open(mets_path)
	if err != nil {
		log.Fatalf("This path does not exist.")
	}
	// We want to read the file
	r := tar.NewReader(file)
	_ = r

	// TODO: We need to loop over the archive to find the mets file.
	// Which should be in the top level of the data directory.

}
