package main

type Premis struct {
	File        string
	Type        []string
	Object      []string
	Outcome     []bool
	Agent       []string
	Format      []string
	EventCount  int
	SuccesCount int
}

// type MetsFile struct {
// 	Name      string
// 	PrType    []*etree.Element
// 	PrObject  []*etree.Element
// 	PrOutcome []*etree.Element
// 	PrAgent   []*etree.Element
// 	PrFormat  []*etree.Element
// }
