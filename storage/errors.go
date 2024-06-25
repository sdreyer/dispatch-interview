package storage

import (
	"auction/auction"
	"fmt"
)

type BidderHasAlreadyBidError struct {
	bidder auction.Bidder
}

func (e *BidderHasAlreadyBidError) Error() string {
	return fmt.Sprintf("bidder %s has already entered a bid", e.bidder)
}

type BidderNotFoundError struct {
	bidder auction.Bidder
}

func (e *BidderNotFoundError) Error() string {
	return fmt.Sprintf("bidder %s not found", e.bidder)
}
