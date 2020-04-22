package rwregister

import (
	"fmt"
	"sync"
)

type Monitor interface {
	RI() int
	RR(i int, v string) (bool, []string)
	WI(v string) int
	WR(i int, v string)
}

type monitorState struct {
	sync.Mutex
	invokeID int
	readers  map[int][]string // map from (read) invokeID to the set of legal values for this reader.
	writers  map[int]string   // map from (write) invokeID to value (hence, the map is the set of all pending values from all writers).
	current  string
}

// NewEventMonitor creates a new monitor.
func NewEventMonitor() Monitor {
	return &monitorState{
		readers: make(map[int][]string),
		writers: make(map[int]string),
	}
}

// RI is ReadInvoke
func (e *monitorState) RI() int {
	e.Lock()
	defer e.Unlock()
	lvals := e.legalValues()
	e.invokeID++
	e.readers[e.invokeID] = lvals
	fmt.Println("RI: Legal values:", lvals, "for the read invoke ID:", e.invokeID)
	return e.invokeID
}

// RR is ReadReturn
func (e *monitorState) RR(i int, v string) (bool, []string) {
	e.Lock()
	defer e.Unlock()
	lvals := e.readers[i]
	delete(e.readers, i)
	fmt.Println("RR: Legal values:", lvals, "for the read return ID:", i)
	return contain(v, lvals), lvals
}

// WI is WriteInvoke
func (e *monitorState) WI(v string) int {
	e.Lock()
	defer e.Unlock()
	e.invokeID++
	e.writers[e.invokeID] = v
	e.updateReaders()
	fmt.Println("WI: Legal values:", e.legalValues(), "for the write invoke ID:", e.invokeID)
	return e.invokeID
}

// WR is WriteReturn
func (e *monitorState) WR(i int, v string) {
	e.Lock()
	defer e.Unlock()
	delete(e.writers, i)
	e.current = v
	fmt.Println("WR: Legal values:", e.legalValues(), "for the write return ID:", i)
}

// legalValues returns the current write and all pending/concurrent writes.
// The caller must hold the lock on the monitor state.
func (e *monitorState) legalValues() []string {
	legalVals := make([]string, 0)
	hasCurrent := false
	for _, val := range e.writers {
		if val == e.current {
			hasCurrent = true
		}
		legalVals = append(legalVals, val)
	}
	if !hasCurrent {
		legalVals = append(legalVals, e.current)
	}
	return legalVals
}

// updateReaders will make sure that all concurrent readers can observe new writes.
// The method is invoked from WI.
// The caller must hold the lock on the monitor state.
func (e *monitorState) updateReaders() {
	lvals := e.legalValues()
	for i := range e.readers {
		for _, lval := range lvals {
			if !contain(lval, e.readers[i]) {
				e.readers[i] = append(e.readers[i], lval)
			}
		}
	}
}

// contain checks if the result from Read quorum call is in the expected value list.
func contain(v string, expectVal []string) bool {
	for _, val := range expectVal {
		if v == val {
			return true
		}
	}
	return len(expectVal) == 0 && len(v) == 0
}
