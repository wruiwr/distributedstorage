package rwregister

import (
	"testing"
)

type Test struct {
	Name       string
	SystemSize uint
	QuorumSize uint
	TestCases  []TestCase
}

type TestCase struct {
	Name           string
	Description    string
	OperationOrder []Operator
}

type Operator interface {
	Operations()
}

type Concurrent []Operation

func (Concurrent) Operations() {}

type Sequential []Operation

func (Sequential) Operations() {}

type Operation struct {
	ID           string
	Name         string
	ArgumentType string
	Value        string
	Replies      ExpectedReplies
}

type Reply struct {
	Type  string
	Value string
}

// ExpectedReplies can hold one or more expected replies.
// Typically, it will be only one possible reply,
// but sometimes it can be more than one possibility.
type ExpectedReplies []Reply

var systemTest = Test{Name: "SystemTest", SystemSize: 3, QuorumSize: 2,
	TestCases: []TestCase{
		{Name: "TestNoConnection", Description: "Test if no connection is handled"},
		{Name: "TestConcurrentReadWriteOps", Description: "Concurrent and Sequential read and write operations",
			OperationOrder: []Operator{
				Concurrent{
					{ID: "A", Name: "Read", ArgumentType: "ReadRequest", Value: "",
						Replies: ExpectedReplies{
							{Type: "Value", Value: "7"},
							{Type: "Value", Value: ""},
						},
					},
					{ID: "B", Name: "Write", ArgumentType: "Value", Value: "7",
						Replies: ExpectedReplies{
							{Type: "WriteResponse", Value: "true"},
						},
					},
				},
				Sequential{
					{ID: "C", Name: "Read", ArgumentType: "ReadRequest", Value: "",
						Replies: ExpectedReplies{
							{Type: "Value", Value: "7"},
						},
					},
					{ID: "D", Name: "Write", ArgumentType: "Value", Value: "8",
						Replies: ExpectedReplies{
							{Type: "WriteResponse", Value: "true"},
						},
					},
				},
				Concurrent{
					{ID: "E", Name: "Write", ArgumentType: "Value", Value: "9",
						Replies: ExpectedReplies{
							{Type: "WriteResponse", Value: "true"},
						},
					},
					{ID: "F", Name: "Read", ArgumentType: "ReadRequest", Value: "",
						Replies: ExpectedReplies{
							{Type: "Value", Value: "8"},
							{Type: "Value", Value: "9"},
						},
					},
				},
				Sequential{
					{ID: "G", Name: "Read", ArgumentType: "ReadRequest", Value: "",
						Replies: ExpectedReplies{
							{Type: "Value", Value: "9"},
						},
					},
				},
			},
		},
	},
}

func TestXSysTest(t *testing.T) {
	t.Logf("%s: (system size=%d, quorum size=%d)", systemTest.Name, systemTest.SystemSize, systemTest.QuorumSize)
	for _, testcase := range systemTest.TestCases {
		t.Logf("%s: %s", testcase.Name, testcase.Description)
		for _, test := range testcase.OperationOrder {
			switch ops := interface{}(test).(type) {
			case Concurrent:
				t.Logf("%v: is concurrent tests", ops)
			case Sequential:
				t.Logf("%v: is sequential tests", ops)
			default:
				t.Errorf("unknown operation type %v", ops)
			}
		}
	}
}
