package id_generator

import "testing"

func WithMemoryIDGenerator() func() IDGenerator {
	return func() IDGenerator {
		return NewMemoryIDGenerator()
	}
}

func Test(t *testing.T) {
	tests := generatorTests{
		generatorFn: WithMemoryIDGenerator(),
		t:           t,
	}
	tests.Run()
}
