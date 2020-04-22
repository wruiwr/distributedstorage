package rwregister

import (
	"fmt"
)

// Run 'go generate' to invoke protoc to compile our protobuf definition.
//go:generate protoc -I=$GOPATH/src/:. --gorums_out=plugins=grpc+gorums:. register.proto

// QuoSpecQ is the quorum specification for the RWRegister quorum algorithm.
type QuoSpecQ struct {
	rq, wq int // quorum sizes for ReadQF and WriteQF
}

// NewQuoSpecQ returns a quorum specification with quorum sizes for
// read and write quorum functions.
func NewQuoSpecQ(rqSize, wqSize int) QuorumSpec {
	return &QuoSpecQ{
		rq: rqSize,
		wq: wqSize,
	}
}

// ReadQF returns nil and false until enough replies have been receive
// so as to satisfy the quorum condition for a read quorum. Hence, once
// enough replies have been received, the highest value and true is returned.
func (qq QuoSpecQ) ReadQF(replies []*Value) (*Value, bool) {
	if len(replies) < qq.rq {
		// fmt.Println("Current quorum size", len(replies), "is not enough for ReadQF")
		// not enough replies yet
		return nil, false
	}

	var highest *Value
	for _, reply := range replies {
		if highest != nil && reply.GetC().GetTimestamp() <= highest.GetC().GetTimestamp() {
			continue
		}
		highest = reply
	}
	fmt.Println("Reply value: ", highest.GetC().GetValue(), " with highest timestamp: ", highest.GetC().GetTimestamp())
	// returns reply with the highest timestamp, or nil if no replies were verified
	return highest, true
}

// WriteQF returns WRITENACK, that is nil and false, until enough replies have been received.
// When enough replies have been received, we return WRITEACK, that is the first reply and true.
func (qq QuoSpecQ) WriteQF(replies []*WriteResponse) (*WriteResponse, bool) {
	if len(replies) < qq.wq {
		// fmt.Println("Current quorum size", len(replies), "is not enough for WriteQF")
		// not enough replies yet
		return nil, false // return (WRITENACK, false)
	}

	return replies[0], true // return (WRITEACK, true)
}
