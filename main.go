package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
)

const (
	ZIP = ".zip"
	TAR = ".tar"
	Z7  = ".7z"
)

func main() {
	// Open the archive
	argAmount := len(os.Args)
	if argAmount < 2 {
		log.Fatalln(`There were not enough arguments.
	 This command requires a path to be given.`)
	} else if argAmount > 3 {
		log.Fatalln(`There are too many arguments.
Only two paths at a time is allowed`)
	}
	// The paths that will be used are the args until the end of the array. We will actually test if they are all valid first.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dst, err := os.MkdirTemp(cwd, "metsFolder")
	defer func() {
		err := os.Remove(dst)
		if err != nil {
			log.Fatal(err)
		}
	}()
	paths := os.Args[1 : len(os.Args)-1]
	for _, path := range paths {
		var tmpMets *os.File
		defer func() {
			if tmpMets != nil {
				err := tmpMets.Close()
				if err != nil {
					log.Fatal("Failed to close tmp")
				}
			}
		}()
		absPath := filepath.Join(cwd, path)
		err := CopyMets(absPath, dst, tmpMets)
		if err != nil {
			log.Fatal(err)
		}
		if tmpMets == nil {
			log.Fatal("No mets file found. Is your mets file not all capitalized?")
		}

		// Create and parse the mets xml file.
		mets := etree.NewDocument()
		if err := mets.ReadFromFile(tmpMets.Name()); err != nil {
			log.Fatalf("Could not parse the mets into an xml file. %v", err)
		}
		root := mets.Root()
		// We need to get all the premis elements from the mets and count them.
		// Use the correct namespace URI in your SelectElement and SelectElements methods.
		// These are placeholders and might need to be adjusted according to your XML.
		// You can find the namespace URI in the XML file, it is the URL specified in the xmlns attribute.

		// Find all premis:eventType elements
		premis := root.FindElements("//premis")
		for i, premisElement := range premis {
			fmt.Printf("\n\n\n%d %v", i, *premisElement)
		}
	}
	if err != nil {
		log.Fatalln("Could not get working directory")
	}

}
