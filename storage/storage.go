package storage

import (
	"auction/currency"
	"auction/id_generator"
)

type Bidder string

type Bid struct {
	Bidder      Bidder
	StartingBid currency.Amount
	MaxBid      currency.Amount
	Increment   currency.Amount
	ID          id_generator.EventID
}

type BidStorer interface {
	SaveBid(bid Bid) error
	GetBid(bidder Bidder) (Bid, error)
	GetAllBids() (map[Bidder]Bid, error)
}
