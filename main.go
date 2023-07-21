package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

const (
	ZIP = ".zip"
	TAR = ".tar"
	Z7  = ".7z"
	XML = ".xml"
)

func main() {
	// Open the archive
	argAmount := len(os.Args)
	if argAmount < 2 {
		log.Fatalln(`There were not enough arguments.

This command requires a path to be given.

USAGE: ardiff <path> <*path...>
			
let * mean optional`)
	} else if argAmount%3 != 0 {
		log.Fatalln(`There most be an even nu
			

Only two paths at a time is allowed`)
	}
	// The paths that will be used are the args until the end of the array. We will actually test if they are all valid first.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dst, err := os.MkdirTemp(cwd, "METS_Data-")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: log a fatal error in every case other than the folder not being empty
	// defer func() {
	// 	err := os.Remove(dst)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	csvF, err := os.Create(filepath.Join(dst, "ComparisonReport.csv"))
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(csvF)
	defer w.Flush()
	header := []string{
		"file", "agent", "outcome", "type", "format"}
	err = w.Write(header)
	if err != nil {
		log.Fatal(err)
	}

	// objectPath := etree.MustCompilePath("//premis:object//*")
	eventPath := etree.MustCompilePath("//premis:event//*")
	// TODO: add this to a configurable setting so that this can be changed
	exclude := map[string]bool{
		"eventIdentifier":                true,
		"eventIdentifierType":            true,
		"eventDateTime":                  true,
		"linkingAgentIdentifierType":     true,
		"linkingAgentIdentifier":         true,
		"linkingAgentIdentifierValue":    false,
		"objectIdentifier":               true,
		"fixity":                         true,
		"objectCharacteristicsExtension": true,
	}

	paths := os.Args[1:len(os.Args)]
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		tmpMets, err := CopyMets(absPath, dst)
		defer func() {
			if tmpMets != nil {
				err := tmpMets.Close()
				if err != nil {
					log.Fatal("Failed to close tmp")
				}
			}
		}()
		if err != nil {
			log.Fatal(err)
		}
		if tmpMets == nil {
			log.Fatal("No mets file found. Is your mets file not all capitalized?")
		}

		metsFile := filepath.Base(tmpMets.Name())
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
		// objs := root.FindElementsPath(objectPath)
		events := root.FindElementsPath(eventPath)
		var etype string
		var outcome string
		var agent string
		var esize string
		var format string
		// var other string
		// Find all premis:eventType elements
		// for _, el := range objs {
		// if _, excluded := exclude[el.Tag]; !excluded {
		// fmt.Print(el.Text(), el.Tag)

		// record := []string{metsFile, agent, el.Text(), etype, format, esize, other}
		// w.Write(record)
		// }
		// }
		for _, el := range events {
			if excluded := exclude[el.Tag]; !excluded {
				switch el.Tag {
				case "eventType":
					etype = el.Text()
					continue
				case "linkingAgentIdentifierValue":
					// TODO: make it configurable which agent identifier is allowed
					if strings.Contains(el.Text(), "Archivematica") || strings.Contains(el.Text(), "a3m") {
						agent = el.Text()
						record := []string{metsFile, agent, outcome, etype, format, esize}
						w.Write(record)
					}
					// case "eventDetail":
					// other = el.Text()
					// continue
				case "eventOutcome":
					// if strings.Contains(el.Text(), "pass") || strings.Contains(el.Text(), "fail") || el.Text() == "" {
					outcome = el.Text()
					continue
					// }
				default:
					// record := []string{metsFile, agent, outcome, etype, format, esize, other}
					// w.Write(record)
				}

			}
		}
	}
	if err != nil {
		log.Fatalln("Could not get working directory")
	}

}
