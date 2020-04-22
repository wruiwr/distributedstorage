package reader

import (
	"encoding/xml"
)

type Test struct {
	XMLName    xml.Name   `xml:"Test"`
	Name       string     `xml:"TestName,attr"`
	SystemSize int        `xml:"SystemSize,attr"`
	QuorumSize int        `xml:"QuorumSize,attr"`
	TestCases  []TestCase `xml:"TestCase"`
}

type TestCase struct {
	XMLName     xml.Name     `xml:"TestCase"`
	Name        string       `xml:"Name,attr"`
	Description string       `xml:"Description,attr"`
	OrderOp     []Operations `xml:"Operations,omitempty"`
}

type Operations struct {
	OpType string      `xml:"Type,attr"`
	Op     []Operation `xml:"Operation"`
}

type Operation struct {
	ID           string  `xml:"ID,attr"`
	Name         string  `xml:"Name,attr"`
	ArgumentType string  `xml:"ArgumentType,attr"`
	Value        string  `xml:"Value,attr,omitempty"`
	Replies      []Reply `xml:"ExpectedReplies>Reply"`
}

type Reply struct {
	Type  string `xml:"Type,attr"`
	Value string `xml:"Value,attr,omitempty"`
}
