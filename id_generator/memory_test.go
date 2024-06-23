package id_generator

import "testing"

func Test(t *testing.T) {
	tests := generatorTests{
		generatorFn: NewMemoryIDGenerator,
		t:           t,
	}
	tests.Run()
}
