package quorumFunTests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	q "github.com/selabhvl/cpnmbt/rwregister"
	r "github.com/selabhvl/cpnmbt/rwregister/reader"
)

var (
	rxmldir     *string
	wxmldir     *string
	readXmlDir  string
	writeXmlDir string
	readTests   r.XMLReadQFTest
	writeTests  r.XMLWriteQFTest
)

func init() {

	rxmldir = flag.String("rdir", "../xml/qftest/Others/newReadQFTests.xml", "the dir to xml files for ReadQF tests")
	wxmldir = flag.String("wdir", "../xml/qftest/Others/newWriteQFTests.xml", "the dir to xml files for WriteQF tests")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	readXmlDir = *rxmldir
	writeXmlDir = *wxmldir
}

// GetReadXml returns test cases in xml for ReadQF
func GetReadXml() r.XMLReadQFTest {
	r.ParseXMLTestCase(readXmlDir, &readTests)
	return readTests
}

// GetWriteXml returns test cases in xml for WriteQF
func GetWriteXml() r.XMLWriteQFTest {
	r.ParseXMLTestCase(writeXmlDir, &writeTests)
	return writeTests
}

// TestReadQF tests ReadQF
func TestReadQF(t *testing.T) {

	quoSpec := NewQuoSpecQ(qSize, qSize)
	quoSpecError := NewQuoSpecQError(qSize, qSize)
	fmt.Println("The test results for", GetReadXml().TestName)

	for _, testcase := range GetReadXml().TestCases {
		fmt.Println("The test case ID:", testcase.CaseID)

		// Convert []*reader.Value to []*gorumsRegister.Value
		values := []*q.Value{}
		for _, re := range testcase.TestValues {
			values = append(values, &q.Value{C: &q.Content{Value: re.TestValue, Timestamp: re.TestTimestamp}})
		}

		fmt.Println("Test correct ReadQF:")
		// Call the correct ReadQF and get testing results
		reply, quorum := quoSpec.ReadQF(values)

		// Compare the testing results of the correct ReadQF against with the test oracle
		if quorum != testcase.ExpectQuorum {
			t.Errorf("In test case %v of the correct ReadQF, got quorum: %t, want quorum: %t", testcase.CaseID, quorum, testcase.ExpectQuorum)
		}

		if reply != nil {
			if reply.C.Value != testcase.ExpectResults.ExpectValue || reply.C.Timestamp != testcase.ExpectResults.ExpectTimestamp {
				t.Errorf("In test case %v of the correct ReadQF, got the reply: %v, want: %v as quorum reply", testcase.CaseID, reply.C, testcase.ExpectResults)
			}
		} else {
			if testcase.ExpectResults != nil {
				t.Errorf("In test case %v of the correct ReadQF, got the reply: %v, want: %v as quorum reply", testcase.CaseID, reply, testcase.ExpectResults)
			}
		}

		fmt.Println("Test the reply with the highest timestamp error of ReadQF:")
		// Test the reply with the highest timestamp error of ReadQF
		reply, quorum = quoSpecError.ReadQF(values)

		// Compare the testing results against with the test oracle
		if quorum != testcase.ExpectQuorum {
			t.Errorf("In test case %v of the reply with the highest timestamp error of ReadQF, got quorum: %t, want quorum: %t", testcase.CaseID, quorum, testcase.ExpectQuorum)
		}

		if reply != nil {
			if reply.C.Value == testcase.ExpectResults.ExpectValue || reply.C.Timestamp == testcase.ExpectResults.ExpectTimestamp {
				t.Errorf("In test case %v of the reply with the highest timestamp error of ReadQF, got the reply: %v, want: %v as quorum reply", testcase.CaseID, reply.C, testcase.ExpectResults)
			}
		} else {
			if testcase.ExpectResults != nil {
				t.Errorf("In test case %v of the reply with the highest timestamp error of ReadQF, got the reply: %v, want: %v as quorum reply", testcase.CaseID, reply, testcase.ExpectResults)
			}
		}

	}

}

// TestWriteQF tests WriteQF
func TestWriteQF(t *testing.T) {

	quoSpec := NewQuoSpecQ(qSize, qSize)
	quoSpecError := NewQuoSpecQError(qSize, qSize)
	fmt.Println("The test results for", GetWriteXml().TestName)

	for _, testcase := range GetWriteXml().TestCases {
		fmt.Println("The test case ID:", testcase.CaseID)

		// Convert []*reader.WriteResponse to []*gorumsRegister.WriteResponse
		writeResponses := []*q.WriteResponse{}
		for _, re := range testcase.TestValues {
			writeResponses = append(writeResponses, &q.WriteResponse{Ack: re.Response})
		}

		fmt.Println("Test correct WriteQF:")
		// Call the correct WriteQF and get testing results
		reply, quorum := quoSpec.WriteQF(writeResponses)

		// Compare the testing results against with the test oracle
		if quorum != testcase.ExpectQuorum {
			t.Errorf("In test case %v of the correct WriteQF, got quorum: %t, want quorum: %t", testcase.CaseID, quorum, testcase.ExpectQuorum)
		}
		if reply != nil { // got the reply
			if reply.Ack != testcase.ExpectResults.Response {
				t.Errorf("In test case %v of the correct WriteQF, got reply: %v, want: %v as quorum reply", testcase.CaseID, reply, testcase.ExpectResults.Response)
			}
		} else { // reply is nil
			fmt.Println("the Response:", testcase.ExpectResults.Response)
			if testcase.ExpectResults.Response {
				t.Errorf("In test case %v of the correct WriteQF, got reply: %v, want: %v as quorum reply", testcase.CaseID, testcase.ExpectResults.Response, !testcase.ExpectResults.Response)
			}
		}

		fmt.Println("Test the number of replies and quorum size error of WriteQF:")
		// Test the number of replies and quorum size error of WriteQF
		reply, quorum = quoSpecError.WriteQF(writeResponses)

		// Compare the testing results against with the test oracle
		if quorum == testcase.ExpectQuorum {
			t.Errorf("In test case %v of the replies and quorum size error of WriteQF, got quorum: %t, want quorum: %t", testcase.CaseID, quorum, testcase.ExpectQuorum)
		}

		if reply != nil {
			t.Logf("The expected results from of writeQF with inject error is: %v", testcase.ExpectResults.Response)
			if reply.Ack == testcase.ExpectResults.Response {
				t.Errorf("In test case %v of the replies and quorum size error of WriteQF, got reply: %v, want: %v as quorum reply", testcase.CaseID, reply.Ack, testcase.ExpectResults.Response)
			}
		} else {
			if testcase.ExpectResults == nil {
				t.Errorf("In test case %v of the replies and quorum size error of WriteQF, got reply: %v, want: %v as expected results", testcase.CaseID, reply, testcase.ExpectResults)
			}
		}
	}
}
