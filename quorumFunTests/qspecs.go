package quorumFunTests

import (
	"fmt"
	// q "github.com/selabhvl/cpnmbt/rwregister"
	q "github.com/wruiwr/distributedstorage"
)

const qSize = 2 // assumed quorum size is 2

// QuoSpecQ is the quorum specification for the Quorum algorithm.
type QuoSpecQ struct {
	rq, wq int // quorum sizes for ReadQF and WriteQF
}

// NewQuoSpecQ has quorum sizes for read and write quorum function as input.
// It returns a quorum spec
func NewQuoSpecQ(rqSize, wqSize int) q.QuorumSpec {

	return &QuoSpecQ{
		rq: rqSize,
		wq: wqSize,
	}
}

// Assumed correct ReadQF:
// Returns nil and false until the supplied replies.
// The method returns the single highest value and true.
func (qq QuoSpecQ) ReadQF(replies []*q.Value) (*q.Value, bool) {
	if len(replies) < qq.rq { // assumed quorum size is 2
		fmt.Println("The number of replies for the correct ReadQF is:", len(replies), ",not enough yet for ReadQF")
		// not enough replies yet
		return nil, false
	}

	// get the reply with the highest timestamp
	var highest *q.Value = &q.Value{C: &q.Content{}}

	for _, reply := range replies {

		if reply.C.Timestamp >= highest.C.Timestamp {
			highest = reply
		}

	}
	fmt.Println("The reply with the highest timestamp for the correct ReadQF is: ", highest)
	// returns reply with the highest timestamp, or nil if no replies were verified
	return highest, true
}

// Assumed correct WriteQF:
// Returns WRITENACK (false or nil) and false until it is possible to check for a quorum.
// If enough replies, we return WRITEACK (true or the first reply) and true.
func (qq QuoSpecQ) WriteQF(replies []*q.WriteResponse) (*q.WriteResponse, bool) {
	if len(replies) < qq.wq {
		fmt.Println("The number of replies for the correct WriteQF is:", len(replies), ",not enough yet for WriteQF")
		// not enough replies yet
		return nil, false // return (WRITENACK, false)
	}

	return replies[0], true // return (WRITEACK, true)
}

// Programming error scenarios for ReadQF and WriteQF
// Error scenario one: number of replies and quorum size error in WriteQF
// Error scenario two: the result of the reply with the highest timestamp error in ReadQF
// QuoSpecQError is the quorum specification for testing QuorumError scenario.
type QuoSpecQError struct {
	rq, wq int // quorum sizes for ReadQF and WriteQF
}

// NewQuoSpecQError has quorum sizes for read and write quorum function as input.
// It returns a quorum spec
func NewQuoSpecQError(rqSize, wqSize int) q.QuorumSpec {

	return &QuoSpecQError{
		rq: rqSize,
		wq: wqSize,
	}
}

// Assumed error ReadQF
// Error scenario for ReadQF: the result of the reply with the highest timestamp error
func (qq QuoSpecQError) ReadQF(replies []*q.Value) (*q.Value, bool) {
	if len(replies) < qq.rq {
		fmt.Println("The number of replies for the Quorum Error ReadQF is:", len(replies), ", not enough yet for ReadQF")
		// not enough replies yet
		return nil, false
	}

	// get the reply with the highest timestamp
	var highest *q.Value = &q.Value{C: &q.Content{}}

	for _, reply := range replies {

		if reply.C.Timestamp < highest.C.Timestamp { // Error
			highest = reply
		}

	}
	fmt.Println("The reply with the highest timestamp for the Quorum Error ReadQF is: ", highest, "when the number of replies is:", len(replies))
	// returns reply with the highest timestamp, or nil if no replies were verified
	return highest, true
}

// Assumed error WriteQF
// Error scenario for WriteQF: number of replies and quorum size error
func (qq QuoSpecQError) WriteQF(replies []*q.WriteResponse) (*q.WriteResponse, bool) {
	if len(replies) >= qq.wq { // Error
		fmt.Println("The number of replies for the correct WriteQF is:", len(replies), ",not enough yet for WriteQF")
		// not enough replies yet
		return nil, false // return (WRITENACK, false)
	}

	return replies[0], true // return (WRITEACK, true)
}
