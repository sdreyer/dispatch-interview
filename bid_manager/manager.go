package bid_manager

import "auction/auction"

type BidManager interface {
	AddBid(bidder, startingBid, maxBid, incrementAmount string) error
	CalculateWinner() (auction.WinningBid, error)
}
