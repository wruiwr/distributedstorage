package rwregister

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

// RegisterTestServer is a basic register server that in addition also can
// signal when a read or write has completed.
type RegisterTestServer interface {
	RegisterServer
	ReadExecuted()
	WriteExecuted()
}

// RegisterServerBasic represents a single State register.
type RegisterServerBasic struct {
	sync.RWMutex
	value Value

	readExecutedChan  chan struct{}
	writeExecutedChan chan struct{}
}

// NewRegisterBasic returns a new basic register server.
func NewRegisterBasic() *RegisterServerBasic {
	return &RegisterServerBasic{
		// Use an appropriate larger buffer size if we construct test
		// scenarios where it's needed.
		value:             Value{C: &Content{Timestamp: time.Now().UnixNano()}},
		writeExecutedChan: make(chan struct{}, 32),
		readExecutedChan:  make(chan struct{}, 32),
	}
}

// Read can handle the reed request from the client.
func (r *RegisterServerBasic) Read(ctx context.Context, rq *ReadRequest) (*Value, error) {
	r.RLock()
	defer r.RUnlock()
	r.readExecutedChan <- struct{}{}

	return &Value{C: r.value.C}, nil
}

// Write can handle the write request from the client.
func (r *RegisterServerBasic) Write(ctx context.Context, s *Value) (*WriteResponse, error) {
	r.Lock()
	defer r.Unlock()
	writeResp := &WriteResponse{}
	if s.C.Timestamp > r.value.C.Timestamp {
		r.value = *s
		writeResp.Ack = true
	}
	r.writeExecutedChan <- struct{}{}

	return writeResp, nil
}

// ReadExecuted returns when r has has completed a read.
func (r *RegisterServerBasic) ReadExecuted() {
	<-r.readExecutedChan
}

// WriteExecuted returns when r has has completed a write.
func (r *RegisterServerBasic) WriteExecuted() {
	<-r.writeExecutedChan
}
