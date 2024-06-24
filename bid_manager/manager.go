package bid_manager

import (
	"auction/currency"
	"auction/storage"
)

type WinningBid struct {
	Bidder storage.Bidder
	Amount currency.Amount
}

// Create a default manager impl
type BidManager interface {
	// Move bids out of storage
	AddBid(bidder, startingBid, maxBid, incrementAmount string) error
	CalculateWinner() (WinningBid, error)
}
