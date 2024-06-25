package bid_manager

import "auction/auction"

type BidManager interface {
	// AddBid creates a bid entry for a person. A person can only enter a single bid entry
	AddBid(bidder, startingBid, maxBid, incrementAmount string) error
	// CalculateWinner returns the winning bid based on the bids that have been added
	CalculateWinner() (auction.WinningBid, error)
}
