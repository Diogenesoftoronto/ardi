package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	z7 "github.com/bodgit/sevenzip"
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
	_ = absPath
	//TODO: There might be an option to create a unified interface for archives (tar, zip, and 7zip)
	ext := filepath.Ext(path)

	dst, err := os.MkdirTemp(cwd, "metsFolder")
	dataPath := filepath.Join(dst, "data")
	var tmpMets *os.File
	defer func() {
		if tmpMets != nil {
			err := tmpMets.Close()
			if err != nil {
				log.Fatal("Failed to close tmp")
			}
		}
	}()
	fmt.Printf("This is the data path %s\n", dataPath)
	// cwd = dataPath
	switch ext {
	case ZIP:
		archive, err := zip.OpenReader(absPath)
		if err != nil {
			log.Fatal(err)
		}
		// You can defer and handle and error by wrapping a function in an anonymous function. This way we can have defer blocks!
		defer func() {
			err := archive.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		for _, f := range archive.File {
			if strings.Contains(f.Name, "mets") {
				tmpMets, err = os.CreateTemp(dst, f.Name)
				if err != nil {
					log.Fatal(err)
				}
				file, err := f.Open()
				if err != nil {
					log.Fatal(err)
				}

				if _, err := io.Copy(tmpMets, file); err != nil {
					log.Fatal("Could not copy the mets")
				}
				break
			}
		}
	case Z7:
		// A bit of code duplication here, I wonder if this really is the best way
		archive, err := z7.OpenReader(absPath)
		if err != nil {
			fmt.Println("here")
			log.Fatal(err)
		}
		defer func() {
			err := archive.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		for _, f := range archive.File {
			if strings.Contains(f.Name, "mets") {
				tmpMets, err = os.CreateTemp(dst, f.Name)
				if err != nil {
					log.Fatal(err)
				}
				file, err := f.Open()
				if err != nil {
					log.Fatal(err)
				}
				if _, err := io.Copy(tmpMets, file); err != nil {
					log.Fatal("Could not copy the mets")
				}
				break
			}
		}
	case TAR:
		r, err := os.Open(absPath)
		if err != nil {
			log.Fatal(err)
		}
		archive := tar.NewReader(r)
		for {
			h, err := archive.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if strings.Contains(h.Name, "mets") {
				tmpMets, err = os.CreateTemp(dst, h.Name)
				if err != nil {
					log.Fatal(err)
				}
				if _, err := io.Copy(tmpMets, archive); err != nil {
					log.Fatal("Could not copy the mets")
				}
				break
			}
		}
	default:
		// In the default case it could be that the folder path we were sent was not compressed or that it is in a format that is not recognized.
		log.Fatal("Currently only compressed files are supported")
	}
	// Now that we are done copying all the mets files to the temp directory we can finally work on them!

	// Create and parse the mets xml file.
	if tmpMets == nil {
		log.Fatal("Fuck this")
	}
	tmpMetsPath := filepath.Join(dst, tmpMets.Name())
	mets := etree.NewDocument()
	if err := mets.ReadFromFile(tmpMetsPath); err != nil {
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
