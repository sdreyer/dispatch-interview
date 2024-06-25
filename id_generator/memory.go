package id_generator

import (
	"sync/atomic"
)

type memoryIDGenerator struct {
	latestID uint32
}

func NewMemoryIDGenerator() IDGenerator {
	return &memoryIDGenerator{
		latestID: 0,
	}
}

// Next uses the atomic library to ensure that concurrent calls are safe and no duplicate EventIDs will be created.
// sync.Mutex could also be used, but this is more efficient for an operation this simple
func (m *memoryIDGenerator) Next() EventID {
	id := atomic.AddUint32(&m.latestID, 1)
	return EventID(id)
}
