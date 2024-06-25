package storage

import (
	"auction/auction"
	"sync"
)

type memoryBidStorage struct {
	bids map[auction.Bidder]auction.Bid
	mtx  *sync.Mutex
}

func NewMemoryBidStorage() BidStorer {
	return &memoryBidStorage{
		bids: map[auction.Bidder]auction.Bid{},
		mtx:  &sync.Mutex{},
	}
}

// SaveBid is a concurrency safe save operation. This is so that if SaveBid and GetBid are called simultaneously
// then it does not result in a concurrent read/write panic and so that GetBid always returns the true set of bids.
func (m memoryBidStorage) SaveBid(bid auction.Bid) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.bids[bid.Bidder]; ok {
		return &BidderHasAlreadyBidError{bidder: bid.Bidder}
	}
	m.bids[bid.Bidder] = bid
	return nil
}

func (m memoryBidStorage) GetBid(bidder auction.Bidder) (auction.Bid, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if bid, ok := m.bids[bidder]; ok {
		return bid, nil
	} else {
		return auction.Bid{}, &BidderNotFoundError{bidder: bid.Bidder}
	}
}

// GetAllBids returns all bids. It currently returns a BidMap instead of a slice to make lookups easier.
// This implementation assumes we are working with a small set of bids. If there were a significant amount of bids
// expected, this would likely work better returning an iterator and using GetBid to lookup specific bidders instead.
func (m memoryBidStorage) GetAllBids() (auction.BidMap, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.bids, nil
}
