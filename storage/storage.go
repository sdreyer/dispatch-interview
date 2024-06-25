package storage

import "auction/auction"

type BidStorer interface {
	SaveBid(bid auction.Bid) error
	GetBid(bidder auction.Bidder) (auction.Bid, error)
	GetAllBids() (auction.BidMap, error)
}
