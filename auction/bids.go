package auction

import (
	"auction/currency"
	"auction/id_generator"
)

type Bidder string
type BidMap map[Bidder]Bid

type Bid struct {
	Bidder      Bidder
	StartingBid currency.Amount
	MaxBid      currency.Amount
	Increment   currency.Amount
	ID          id_generator.EventID
}

type WinningBid struct {
	Bidder Bidder
	Amount currency.Amount
}
