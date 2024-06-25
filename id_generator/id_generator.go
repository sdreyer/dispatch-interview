package id_generator

type EventID uint32

type IDGenerator interface {
	// Next returns the next EventID. This is a concurrency safe operation and any other implementations
	// also need to be concurrency safe.
	Next() EventID
}
