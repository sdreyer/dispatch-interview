package bid_manager

import "fmt"

type EmptyBidListError struct {
}

func (e *EmptyBidListError) Error() string {
	return fmt.Sprintf("cannot calculate bids. no bids have been entered")
}

type InvalidBidError struct {
	message string
}

func (e *InvalidBidError) Error() string {
	return e.message
}
