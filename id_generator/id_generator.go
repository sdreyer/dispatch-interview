package id_generator

type EventID uint32

type IDGenerator interface {
	Next() EventID
}
