package main

import (
	"encoding/csv"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/beevik/etree"
	"github.com/charmbracelet/log"
)

func main() {
	// Open the archive
	argAmount := len(os.Args)
	if argAmount < 3 {
		log.Fatal(`There were not enough arguments.

This command requires a path to be given.

USAGE: ardi <path> <path...>`)
	} else if argAmount%3 != 0 {
		log.Fatal(`There must be an even number of arguments given
			
Only two paths at a time is allowed`)
	}
	// The paths that will be used are the args until the end of the
	// array. We will actually test if they are all valid first.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dst, err := os.MkdirTemp(cwd, "Mets_Data-")
	if err != nil {
		log.Fatal(err)
	}

	csvF, err := os.Create(filepath.Join(dst, "Comparison_Report.csv"))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("The comparison report is located in %s", dst)
	w := csv.NewWriter(csvF)
	defer w.Flush()
	// example result of the csv
	// file_1,file_2,events_1,events_2,agent_1,agent_2,eventCount_1,eventCount_2,successCount_1,successCount_2
	// mets-2349.xml,mets-3453.xml,{[1:{"id":"<uuid>","format": "excel", "type": "creation", "outcome": "pass"}]},{[1:{"id":"<uuid>","format": "excel", "type": "creation", "outcome": "pass"}]},Archivematica,a3m,1,1,1,1
	header := []string{
		"file_1", "file_2", "event_diff", "agent_1", "agent_2", "eventCount_1", " eventCount_2", " successCount_1", " succussCount_2"}
	err = w.Write(header)
	if err != nil {
		log.Fatal(err)
	}

	paths := os.Args[1:len(os.Args)]
	data := make([]MetsData, len(paths))
	for i, path := range paths {
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
		data[i].File = metsFile
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
		// This is all the amd sections for each mets.
		amdSecs := root.FindElementsPath(amdSecPath)
		// We need to diff the events somehow, I am considering that we just use the difftool for now and then add the diff to the csv.

		// At some point we could decide to go through all the events figureout how
		// many there are before handling them and assigning an slice with that length
		// in order to save memory but I don't think that is worth it.

		eventTotal := len(root.FindElementsPath(eventAmountPath))
		data[i].EventCount = eventTotal
		for _, sec := range amdSecs {
			data[i].handleEvents(sec)
		}
		// Let's create two json files for each of the mets, call them the corresponding name of the mets files.
		var f1, f2 *os.File
		if (i+1)%2 == 0 {
			f1, err = os.Create(data[i-1].File + ".json")
			if err != nil {
				panic(err)
			}
			f2, err = os.Create(data[i].File + ".json")
			if err != nil {
				panic(err)
			}
		} else {
			continue
		}
		defer f1.Close()
		defer f2.Close()
		json1, err := serializeEvents(data[i-1].Events)
		if err != nil {
			panic(err)
		}
		json2, err := serializeEvents(data[i].Events)
		if err != nil {
			panic(err)
		}
		_, err = f1.Write(json1)
		if err != nil {
			panic(err)
		}
		_, err = f2.Write(json2)
		if err != nil {
			panic(err)
		}
		diffCmd := exec.Command("diff", "-u", f1.Name(), f2.Name())
		out, _ := diffCmd.Output()
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				log.Warnf("Ardi found differences between %s and %s", f1.Name(), f2.Name())
			}
		}
		if string(out) == "" {
			log.Info("Ardi found no differences. The Premis events are Identical.")
			os.Exit(0)
		}
		ec1 := strconv.Itoa(data[i-1].EventCount)
		ec2 := strconv.Itoa(data[i].EventCount)
		es1 := strconv.Itoa(data[i-1].SuccesCount)
		es2 := strconv.Itoa(data[i].SuccesCount)
		w.Write([]string{
			data[i-1].File,
			data[i].File,
			string(out),
			data[i-1].Agent,
			data[i].Agent,
			ec1, ec2, es1, es2,
		})

	}
	if err != nil {
		log.Fatal("Could not get working directory")
	}

}
