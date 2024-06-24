package bid_manager

import (
	"auction/currency"
	"auction/id_generator"
	"auction/storage"
	"errors"
	"fmt"
)

type bidState map[storage.Bidder]currency.Amount

type MemoryBidManager struct {
	idGenerator id_generator.IDGenerator
	storage     storage.BidStorer
}

// Change this so that we pass the generator and storage in
func NewMemoryBidManager() (BidManager, error) {
	return &MemoryBidManager{
		// These should probably have an option to throw an error on init
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}, nil
}

func (m MemoryBidManager) AddBid(bidder, startingBid, maxBid, incrementAmount string) error {
	start, err := currency.ParseAmount(startingBid)
	if err != nil {
		return errors.Join(errors.New("failed to parse starting bid"), err)
	}

	//rename maxb
	maxB, err := currency.ParseAmount(maxBid)
	if err != nil {
		return errors.Join(errors.New("failed to parse max bid"), err)
	}

	increment, err := currency.ParseAmount(incrementAmount)
	if err != nil {
		return errors.Join(errors.New("failed to parse increment amount"), err)
	}

	bid := storage.Bid{
		Bidder:      storage.Bidder(bidder),
		StartingBid: start,
		MaxBid:      maxB,
		Increment:   increment,
		ID:          m.idGenerator.Next(),
	}

	err = m.storage.SaveBid(bid)
	if err != nil {
		return errors.Join(errors.New("failed to save bid"), err)
	}
	return nil
}

func (m MemoryBidManager) CalculateWinner() (WinningBid, error) {
	var currentWinner WinningBid

	bids, err := m.storage.GetAllBids()
	if err != nil {
		return WinningBid{}, errors.Join(errors.New("failed to fetch bids"), err)
	}

	state := m.initializeCalculation(bids)

	complete := false
	for !complete {
		fmt.Println("CHECKING")
		state = m.calculateBids(bids, state, currentWinner)
		currentWinner = m.currentWinner(bids, state, currentWinner)
		complete = m.isFinished(bids, state, currentWinner)
		fmt.Println(state)
		fmt.Println(currentWinner)
		fmt.Println(complete)
	}
	return currentWinner, nil
}

func (m MemoryBidManager) initializeCalculation(bids map[storage.Bidder]storage.Bid) bidState {
	state := bidState{}

	for bidder, bid := range bids {
		state[bidder] = bid.StartingBid
	}
	return state
}

// Refactor this
func (m MemoryBidManager) calculateBids(bids map[storage.Bidder]storage.Bid, state bidState, currentWinner WinningBid) bidState {
	newState := bidState{}
	for bidder, amount := range state {
		newAmount := amount
		bid := bids[bidder]
		for (currentWinner.Amount.Greater(newAmount) || currentWinner.Amount.Equals(newAmount)) && (bid.MaxBid.Sub(newAmount).Greater(bid.Increment) || bid.MaxBid.Sub(newAmount).Equals(bid.Increment)) {
			newAmount = newAmount.Add(bid.Increment)
		}
		newState[bidder] = newAmount
	}
	return newState
}

// Add test case for ties
func (m MemoryBidManager) currentWinner(bids map[storage.Bidder]storage.Bid, state bidState, currentWinner WinningBid) WinningBid {
	highestBidder := currentWinner.Bidder
	for bidder, amount := range state {
		if amount.Greater(state[highestBidder]) {
			highestBidder = bidder
		} else if amount.Equals(state[highestBidder]) {
			if bids[bidder].ID < bids[highestBidder].ID {
				highestBidder = bidder
			}
		}
	}
	return WinningBid{
		highestBidder,
		state[highestBidder],
	}
}

func (m MemoryBidManager) isFinished(bids map[storage.Bidder]storage.Bid, state bidState, currentWinner WinningBid) bool {
	complete := true
	for bidder, amount := range state {
		bid := bids[bidder]
		if (bid.MaxBid.Sub(amount).Greater(bid.Increment) || bid.MaxBid.Sub(amount).Equals(bid.Increment)) && bidder != currentWinner.Bidder {
			complete = false
		}
	}
	return complete
}

// If users bid is less than the winning bid AND users max bid minus increment is greater than increment, then add increment
// until users current bid is greater than current winning bid OR until users max bid minus current bid is less than increment
