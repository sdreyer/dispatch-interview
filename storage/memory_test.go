package storage

import "testing"

func WithMemoryBidStorage() func() BidStorer {
	return func() BidStorer {
		return NewMemoryBidStorage()
	}
}

func Test(t *testing.T) {
	tests := storageTests{
		storeFn: WithMemoryBidStorage(),
		t:       t,
	}
	tests.Run()
}
