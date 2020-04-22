package reader

import (
	"encoding/xml"
	"io/ioutil"
)

// XMLSysTest is an exported type for System tests
type XMLSysTest struct {
	XMLName   xml.Name          `xml:"Test"`
	TestName  string            `xml:"TestName,attr"`
	TestCases []XMLSysTestCases `xml:"TestCase"`
}

type XMLSysTestCases struct {
	XMLName  xml.Name             `xml:"TestCase"`
	CaseID   string               `xml:"CaseID,attr"`
	Routines []XMLSysTestRoutines `xml:"Routine"`
}

type XMLSysTestRoutines struct {
	XMLName       xml.Name             `xml:"Routine"`
	RoutineID     string               `xml:"RoutineID,attr"`
	OperationName string               `xml:"OperationName,attr"`
	RoutineValues []string             `xml:"OperationValues>Value"`
	Routines      []XMLSysTestRoutines `xml:"Routine"`
}

// XMLReadQFTest is an exported type for ReadQF tests
type XMLReadQFTest struct {
	XMLName   xml.Name             `xml:"Test"`
	TestName  string               `xml:"TestName,attr"`
	TestCases []XMLReadQFTestCases `xml:"TestCase"`
}

type XMLReadQFTestCases struct {
	XMLName       xml.Name                `xml:"TestCase"`
	CaseID        string                  `xml:"CaseID,attr"`
	TestValues    []*XMLReadQFTestValues  `xml:"TestValues>Content"`
	ExpectResults *XMLReadQFExpectResults `xml:"ExpectResults"`
	ExpectQuorum  bool                    `xml:"ExpectQuorum"`
}

type XMLReadQFTestValues struct {
	XMLName       xml.Name `xml:"Content"`
	TestValue     string   `xml:"Value"`
	TestTimestamp int64    `xml:"Timestamp"`
}

type XMLReadQFExpectResults struct {
	XMLName         xml.Name `xml:"ExpectResults"`
	ExpectValue     string   `xml:"Value"`
	ExpectTimestamp int64    `xml:"Timestamp"`
}

// XMLWriteQFTest is an exported type for WriteQF tests
type XMLWriteQFTest struct {
	XMLName   xml.Name              `xml:"Test"`
	TestName  string                `xml:"TestName,attr"`
	TestCases []XMLWriteQFTestCases `xml:"TestCase"`
}

type XMLWriteQFTestCases struct {
	XMLName       xml.Name                 `xml:"TestCase"`
	CaseID        string                   `xml:"CaseID,attr"`
	TestValues    []*XMLWriteQFTestValues  `xml:"TestValues>Content"`
	ExpectResults *XMLWriteQFExpectResults `xml:"ExpectResults"`
	ExpectQuorum  bool                     `xml:"ExpectQuorum"`
}

type XMLWriteQFTestValues struct {
	XMLName  xml.Name `xml:"Content"`
	Response bool     `xml:"Response"`
}

type XMLWriteQFExpectResults struct {
	XMLName  xml.Name `xml:"ExpectResults"`
	Response bool     `xml:"Response"`
}

func ParseXMLTestCase(file string, xmlTestCaseType interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, &xmlTestCaseType)
}
