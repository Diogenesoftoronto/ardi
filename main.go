package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/beevik/etree"
	"github.com/charmbracelet/log"
	"golang.org/x/exp/slices"
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
	csvDoc := make([][]string, 0)
	defer w.Flush()
	// example result of the csv
	// file_1,file_2,events_1,events_2,agent_1,agent_2,eventCount_1,eventCount_2,successCount_1,successCount_2
	// mets-2349.xml,mets-3453.xml,{[1:{"id":"<uuid>","format": "excel", "type": "creation", "outcome": "pass"}]},{[1:{"id":"<uuid>","format": "excel", "type": "creation", "outcome": "pass"}]},Archivematica,a3m,1,1,1,1
	header := []string{
		"transfer", "file_1", "file_2", "event_diff", "agent_1", "agent_2", "eventCount_1", " eventCount_2", "successCount_1", "successCount_2", "eventTypeCount_1", "eventTypeCount_2"}
	csvDoc = append(csvDoc, header)
	defer func(doc *[][]string) {
		w.WriteAll(*doc)
	}(&csvDoc)

	// todo: consider handling the errors if the writer fails at some point
	// err = w.Write(header)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	paths := os.Args[1:len(os.Args)]
	data := make([]FileData, len(paths))
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

		// Get the name of the tranfer for the mets file
		transfer := root.FindElementPath(transferNamePath).Text()
		// Retrieve only the non slash characters
		re := regexp.MustCompile("/")
		transfer = string(re.ReplaceAll([]byte(transfer), []byte("")))
		data[i].Transfer = transfer
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
			f1, err = os.Create(filepath.Join(dst, data[i-1].File+".json"))
			if err != nil {
				panic(err)
			}
			f2, err = os.Create(filepath.Join(dst, data[i].File+".json"))
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
		csvDoc = append(csvDoc, []string{
			data[i-1].File,
			data[i].File,
			string(""),
			data[i-1].Agent,
			data[i].Agent,
			ec1, ec2, es1, es2,
		})
		// After writing for the total counts for each mets file,
		// we need to write the results for the other premi objects.
		// We will need to loop through each object in the mets and
		// output the results in the fields of the csv.
		dd1, dd2 := convertAllEvents(data[i-1].Events,
			data[i-1].Agent), convertAllEvents(data[i].Events,
			data[i].Agent)
		// log.Infof("%v \n %v", dd1, dd2)

		// let's go through all unique keys for dd1, and dd2.
		// if the key in one does not exist in the other we will need to do a check.
		entries := [][]any{}
		for k, v := range dd1 {
			_, ok := dd2[k]
			if !ok {
				dd2[k] = PremisData{}
			}
			entries = append(entries, []any{k, v})
		}
		for k, v := range dd2 {
			_, ok := dd1[k]
			if !ok {
				dd1[k] = PremisData{}
				entries = append(entries, []any{k, v})
			}
		}
		eTypes := []string{}
		for _, entr := range entries {
			// pre := fmt.Sprintf("+++%s\n---%s\n\n", filepath.Base(f1.Name()), filepath.Base(f2.Name()))
			diff := ""
			k, _ := entr[0].(string)
			v, _ := entr[1].(PremisData)
			// We can think of each array of events as a set. If we
			// think in this way we only need to check if an item
			// in one set is contained in the other. This requires
			// going through each set. But on the second set we don't
			// need to check if it contains every item. Instead if
			// we can get the id of the events then we can simply
			// look at the length and determine which items are
			// not contained in the first set.
			// arg1, arg2 := v.Events[id], dd2[k].Events
			etm := make(map[string]int)
			for _, e := range dd1[k].Events {
				etm[e]++
				if !slices.Contains(dd2[k].Events, e) {
					diff += fmt.Sprintf("+++\t%s\n", e)
				}
				if !slices.Contains(eTypes, e) {
					eTypes = append(eTypes, e)
					// Once we come across a unique event we will record how
					// many such events occured in total.

				}

			}
			etm2 := make(map[string]int)
			for _, e := range dd2[k].Events {
				etm2[e]++
				if !slices.Contains(v.Events, e) {
					diff += fmt.Sprintf("---\t%s\n", e)
				}
				if !slices.Contains(eTypes, e) {
					eTypes = append(eTypes, e)
				}
			}

			jsd, err := json.MarshalIndent(etm, "", "\t")
			if err != nil {
				log.Warn(err)
			}

			jsd2, err := json.MarshalIndent(etm2, "", "\t")
			if err != nil {
				log.Warn(err)
			}
			csvDoc = append(csvDoc, []string{
				k,
				k,
				diff,
				dd1[k].Agent,
				dd2[k].Agent,
				strconv.Itoa(dd1[k].EventCount),
				strconv.Itoa(dd2[k].EventCount),
				strconv.Itoa(dd1[k].SuccessCount),
				strconv.Itoa(dd2[k].SuccessCount),
				string(jsd),
				string(jsd2),
			})
		}

		// csvDoc[0] = append(csvDoc[0], eTypes...)
	}

	// log.Info(csvDoc)
	if err != nil {
		log.Fatal("Could not get working directory")
	}

}
