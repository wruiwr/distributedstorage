package rwregister

import "testing"

var tests = []struct {
	inStr    string
	inSlice  []string
	expected bool
}{
	{"f", nil, false},
	{"f", []string{}, false},
	{"", nil, true},        //TODO should this be true?
	{"", []string{}, true}, //TODO should this be true??
	{"", []string{"", "a", "b", "c"}, true},
	{"a", []string{"a", "b", "c"}, true},
	{"cf", []string{"a", "b", "c", "d", "e", "cf"}, true},
	{"ck", []string{"a", "b", "c", "cafe", "bafe", "ck"}, true},
	{"bahkjha", []string{"a", "b", "c", "fjehf", "fefef", "fekjhfejkh", "bahkjha"}, true},
	{"hihoh", []string{"a", "b", "c", "fjehf", "fefef", "fekjhfejkh", "bahkjha"}, false},
	{"hihoh", []string{"a", "b", "c", "fjehf", "fefef", "fekjhfejkh", "bahkjha", "jefhejh", "h"}, false},
	{"hihoh", []string{"a", "b", "c", "fjehf", "fefef", "fekjhfejkh", "bahkjha", "fekjfhe", "fekjfhejkh"}, false},
}

func TestContain(t *testing.T) {
	for i, test := range tests {
		if !contain(test.inStr, test.inSlice) == test.expected {
			t.Errorf("%d: got %t, expected %t for (%s ∈ %v)", i, !test.expected, test.expected, test.inStr, test.inSlice)
		}
	}
}

// run with go test -bench . -run=BenchmarkContain

func BenchmarkContain(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			contain(test.inStr, test.inSlice)
		}
	}
}
