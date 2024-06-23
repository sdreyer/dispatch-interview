package id_generator

import "sync"

type Memory struct {
	latestID EventID
	mtx      *sync.Mutex
}

func NewMemoryIDGenerator() IDGenerator {
	return &Memory{
		latestID: 0,
		mtx:      &sync.Mutex{},
	}
}

func (m *Memory) Next() EventID {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.latestID = m.latestID + 1
	return m.latestID
}
