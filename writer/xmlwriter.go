package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Test struct {
	XMLName    xml.Name   `xml:"Test"`
	Name       string     `xml:"TestName,attr"`
	SystemSize uint       `xml:"SystemSize,attr"`
	QuorumSize uint       `xml:"QuorumSize,attr"`
	TestCases  []TestCase `xml:"TestCase"`
}

type TestCase struct {
	XMLName     xml.Name     `xml:"TestCase"`
	Name        string       `xml:"Name,attr"`
	Description string       `xml:"Description,attr"`
	OrderOp     []Operations `xml:"Operations,omitempty"`
}

type Operations struct {
	OpType string       `xml:"Type,attr"`
	Op     *[]Operation `xml:"Operation"`
}

type Operation struct {
	ID           string   `xml:"ID,attr"`
	Name         string   `xml:"Name,attr"`
	ArgumentType string   `xml:"ArgumentType,attr"`
	Value        string   `xml:"Value,attr,omitempty"`
	Replies      *[]Reply `xml:"ExpectedReplies>Reply"`
}

type Reply struct {
	Type  string `xml:"Type,attr"`
	Value string `xml:"Value,attr,omitempty"`
}

func main() {
	xmlWriter("./tests.xml")
}

// xmlWriter can write test cases for ReadQF, WriteQF or System tests into xml files.
func xmlWriter(dir string) {
	output, err := xml.MarshalIndent(result, " ", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	testFilePath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(testFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	_, err = f.Write(output)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Wrote XML output to:\n", testFilePath)
}

var result = &Test{
	Name:       "SystemTest",
	SystemSize: 3,
	QuorumSize: 2,
	TestCases: []TestCase{
		{
			Name:        "TestNoConnection",
			Description: "Test if no connection is handled",
		},
		{
			Name:        "TestConcurrentReadWriteOps",
			Description: "Concurrent and Sequential read and write operations",
			OrderOp: []Operations{
				{
					OpType: "Concurrent",
					Op: &[]Operation{
						{
							ID:           "A",
							Name:         "DoReadCall",
							ArgumentType: "ReadRequest",
							Value:        "",
							Replies: &[]Reply{
								{Type: "Value", Value: "7"},
								{Type: "Value", Value: ""},
							},
						},

						{
							ID:           "B",
							Name:         "DoWriteCall",
							ArgumentType: "Value",
							Value:        "7",
							Replies: &[]Reply{
								{Type: "WriteResponse", Value: "true"},
							},
						},
					},
				},
				{
					OpType: "Sequential",
					Op: &[]Operation{
						{
							ID:           "C",
							Name:         "DoReadCall",
							ArgumentType: "ReadRequest",
							Value:        "",
							Replies: &[]Reply{
								{Type: "Value", Value: "7"},
							},
						},
						{
							ID:           "C",
							Name:         "DoWriteCall",
							ArgumentType: "Value",
							Value:        "8",
							Replies: &[]Reply{
								{Type: "WriteResponse", Value: "true"},
							},
						},
						{
							ID:           "C",
							Name:         "ServerFailure",
							ArgumentType: "Value",
							Value:        "1",
						},
					},
				},
				{
					OpType: "Concurrent",
					Op: &[]Operation{
						{
							ID:           "A",
							Name:         "DoReadCall",
							ArgumentType: "ReadRequest",
							Value:        "",
							Replies: &[]Reply{
								{Type: "Value", Value: "9"},
								{Type: "Value", Value: "8"},
							},
						},

						{
							ID:           "B",
							Name:         "DoWriteCall",
							ArgumentType: "Value",
							Value:        "9",
							Replies: &[]Reply{
								{Type: "WriteResponse", Value: "true"},
							},
						},
					},
				},
				{
					OpType: "Sequential",
					Op: &[]Operation{
						{
							ID:           "C",
							Name:         "DoReadCall",
							ArgumentType: "ReadRequest",
							Value:        "",
							Replies: &[]Reply{
								{Type: "Value", Value: "9"},
							},
						},
					},
				},
				{
					OpType: "Concurrent",
					Op: &[]Operation{
						{
							ID:           "A",
							Name:         "DoWriteCall",
							ArgumentType: "Value",
							Value:        "10",
							Replies: &[]Reply{
								{Type: "WriteResponse", Value: "true"},
							},
						},

						{
							ID:           "B",
							Name:         "DoWriteCall",
							ArgumentType: "Value",
							Value:        "11",
							Replies: &[]Reply{
								{Type: "WriteResponse", Value: "true"},
							},
						},
					},
				},
				{
					OpType: "Sequential",
					Op: &[]Operation{
						{
							ID:           "C",
							Name:         "DoReadCall",
							ArgumentType: "ReadRequest",
							Value:        "",
							Replies: &[]Reply{
								{Type: "Value", Value: "10"},
								{Type: "Value", Value: "11"},
							},
						},
					},
				},
			},
		},
	},
}
