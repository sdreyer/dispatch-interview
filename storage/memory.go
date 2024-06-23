package storage

import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	bids map[Bidder]Bid
	mtx  *sync.Mutex
}

func NewMemoryBidStorage() BidStorer {
	return &MemoryStorage{
		bids: map[Bidder]Bid{},
		mtx:  &sync.Mutex{},
	}
}

func (m MemoryStorage) SaveBid(bid Bid) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.bids[bid.Bidder]; ok {
		return fmt.Errorf("bidder %s has already entered a bid", bid.Bidder)
	}
	m.bids[bid.Bidder] = bid
	return nil
}

func (m MemoryStorage) GetBid(bidder Bidder) (Bid, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if bid, ok := m.bids[bidder]; ok {
		return bid, nil
	} else {
		// Should make this nullable
		return Bid{}, fmt.Errorf("cannot find bidder %s", bidder)
	}
}

// Maybe turn this into an iterator later
func (m MemoryStorage) GetAllBids() ([]Bid, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	var bids []Bid
	for _, bid := range m.bids {
		bids = append(bids, bid)
	}
	return bids, nil
}
