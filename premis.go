package main

type MetsData struct {
	Id          int
	File        string
	Events      []Event
	Agent       string // e.g. Archivematica, a3m
	EventCount  int    // e.g. len(events)
	SuccesCount int    // e.g. the amount of event.outcome that are positive or pass
}

type Event struct {
	Id            string `json:"id"` //uuid type of field taken from the mets
	OutcomeDetail string `json:detail`
	Type          string `json:"type"`    //event type e.g. fixity check, creation
	ObjectName    string `json`           //premisObjectOrginalName
	Outcome       bool   `json:"outcome"` //can be empty, but this one is weird e.g. pass, Positive, etc.
}

// file_1,file_2,events_1,events_2,agent_1,agent_2,eventCount_1,eventCount_2,successCount_1,successCount_2
// mets-2349.xml,mets-3453.xml,{[1:{"id":"<uuid>","format": "excel", "type": "creation", "outcome": "pass"}]},{[1:{"id":"<uuid>","format": "excel", "type": "creation", "outcome": "pass"}]},Archivematica,a3m,1,1,1,1
