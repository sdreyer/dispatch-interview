package storage

import (
	"auction/auction"
	"sync"
)

type MemoryStorage struct {
	bids map[auction.Bidder]auction.Bid
	mtx  *sync.Mutex
}

func NewMemoryBidStorage() BidStorer {
	return &MemoryStorage{
		bids: map[auction.Bidder]auction.Bid{},
		mtx:  &sync.Mutex{},
	}
}

func (m MemoryStorage) SaveBid(bid auction.Bid) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.bids[bid.Bidder]; ok {
		return &BidderHasAlreadyBidError{bidder: bid.Bidder}
	}
	m.bids[bid.Bidder] = bid
	return nil
}

func (m MemoryStorage) GetBid(bidder auction.Bidder) (auction.Bid, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if bid, ok := m.bids[bidder]; ok {
		return bid, nil
	} else {
		return auction.Bid{}, &BidderNotFoundError{bidder: bid.Bidder}
	}
}

// Maybe turn this into an iterator later
func (m MemoryStorage) GetAllBids() (auction.BidMap, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.bids, nil
}
