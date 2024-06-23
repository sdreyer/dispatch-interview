package id_generator

import (
	"sync"
	"testing"
)

type generatorTests struct {
	// Make this more generic later
	generatorFn func() IDGenerator
	t           *testing.T
}

func (g *generatorTests) Run() {
	tests := map[string]func(t *testing.T, generator IDGenerator){
		"Test Start Value":          testStartValue,
		"Test Serial Increment":     testSerialIncrement,
		"Test Concurrent Increment": testConcurrentIncrement,
	}
	for name, test := range tests {
		g.t.Run(name, func(t *testing.T) {
			test(t, g.generatorFn())
		})
	}
}

func testStartValue(t *testing.T, generator IDGenerator) {
	initialValue := generator.Next()
	if initialValue != 1 {
		t.Fatalf("Expected initial value to be 1, got: %d", initialValue)
	}
}

func testSerialIncrement(t *testing.T, generator IDGenerator) {
	for i := 1; i <= 100; i++ {
		val := generator.Next()
		if EventID(i) != val {
			t.Fatalf("Expected value %d, got %d", i, val)
		}
	}
}

func testConcurrentIncrement(t *testing.T, generator IDGenerator) {
	numRoutines := 1000
	values := make(chan EventID, numRoutines)
	wg := &sync.WaitGroup{}
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			values <- generator.Next()
		}()
	}
	wg.Wait()

	if len(values) != numRoutines {
		t.Fatalf("Expected %d value, got: %d", numRoutines, len(values))
	}

	idSet := map[EventID]struct{}{}
	for len(values) != 0 {
		value := <-values
		if _, ok := idSet[value]; ok {
			t.Fatalf("Duplicate value found: %d", value)
		}
		idSet[value] = struct{}{}
	}
}
