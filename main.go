package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	z7 "github.com/bodgit/sevenzip"
)

func main() {
	// Open the archive
	argAmount := len(os.Args)
	if argAmount < 2 {
		log.Fatalln(`There were not enough arguments.
	 This command requires a path to be given.`)
	} else if argAmount > 2 {
		// TODO: Allow multiple paths
		log.Fatalln(`There are too many arguments.
Only a single path at a time is allowed`)
	}

	path := os.Args[1]
	fmt.Printf("here are some path:\n%s\n", path)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get working directory")
	}
	// Need to look for a file that includes a file with a mets.
	// The mets will have a uuid: mets.<uuid>.xml
	// If this is a tar file or zip or 7z we need to extract and read from that directory before going to the data directory
	absPath := filepath.Join(cwd, path)
	ext := filepath.Ext(path)
	switch ext {
	case ".zip":
		fmt.Print("Using a zip")
		fallthrough
	case ".tar":
		fmt.Print("or using a tar")
		fallthrough
	case ".7z":

		archive, err := z7.OpenReader(absPath)
		if err != nil {
			log.Fatalf("Could not open archive %s. got:\n %v", absPath, err)
		}
		for _, f := range archive.File {
			fmt.Printf("%v\n", f)
		}

	default:
		fmt.Print("ahhhhhhhhhhhhh")
	}
	dataPath := filepath.Join(cwd, path, "data")
	// if err != nil {
	// 	log.Fatalf("The relative path could not be found. path: %s\ndatapath: %s\ncwd: %s\n", path, dataPath, cwd)
	// }
	fmt.Printf("This is the path %s\n", dataPath)
	cwd = dataPath
	dir, err := os.ReadDir(cwd)
	if err != nil {
		log.Fatalln("This path does not exist.")
	}
	var entry os.DirEntry

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
		fmt.Printf("%d %v", i, *premisEvent)
	}
	// TODO: We need to loop over the archive to find the mets file.
	// Which should be in the top level of the data directory.

}
