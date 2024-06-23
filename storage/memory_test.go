package storage

import "testing"

func Test(t *testing.T) {
	tests := storageTests{
		storeFn: NewMemoryBidStorage,
		t:       t,
	}
	tests.Run()
}
