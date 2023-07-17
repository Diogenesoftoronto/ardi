package main

import (
	"archive/tar"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Open the archive
	var path string

	mets_path := filepath.Join(path, filepath.Join("data", "mets.xml"))
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("This path does not exist.")
	}
	// We want to read the file
	r := tar.NewReader(file)
	// We need to loop over the archive to find the mets file.
	// Which should be in the top level of the data directory.

}
