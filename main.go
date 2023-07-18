package main

import (
	"fmt"
	"github.com/beevik/etree"
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
	fmt.Printf("here are some args:\n%s\n", args)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory")
	}
	// Need to look for a file that includes a file with a mets.
	// The mets will have a uuid: mets.<uuid>.xml
	dataPath := filepath.Join(cwd, os.path, "data")
	cwd = dataPath
	dir, err := os.ReadDir(dataPath)
	if err != nil {:
		log.Fatalf("This path does not exist.")
	}
	var entry os.DirEntry:w
	q
	for _, dEntry := range dir {
		if strings.Contains(dEntry.Name(), "mets") && !dEntry.IsDir() {
			entry = dEntry
			break
		}
	}

	// Create and parse the mets xml file.
	mets := etree.NewDocument()
	if err := mets.ReadFromFile(entry.Name()); err != nil {
		log.Fatalf("Could not parse the mets into an xml file. %v", err)
	}
	root := mets.SelectElement("mets:mets")
	// We need to get all the premis elements from the mets and count them.
	for i, premisEvent := range root.SelectElements("premis:eventType") {
		fmt.Printf("%d %s", i, &premisEvent)
	}
	// TODO: We need to loop over the archive to find the mets file.
	// Which should be in the top level of the data directory.

}
