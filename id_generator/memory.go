package id_generator

import (
	"sync/atomic"
)

type Memory struct {
	latestID uint32
}

func NewMemoryIDGenerator() IDGenerator {
	return &Memory{
		latestID: 0,
	}
}

func (m *Memory) Next() EventID {
	id := atomic.AddUint32(&m.latestID, 1)
	return EventID(id)
}
