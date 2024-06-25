package bid_manager

import (
	"auction/auction"
	"auction/currency"
	"auction/id_generator"
	"auction/storage"
	"errors"
	"fmt"
)

type bidState map[auction.Bidder]currency.Amount

type DefaultBidManager struct {
	idGenerator id_generator.IDGenerator
	storage     storage.BidStorer
}

func NewDefaultBidManager(idGenerator id_generator.IDGenerator, store storage.BidStorer) (BidManager, error) {
	return &DefaultBidManager{
		idGenerator: idGenerator,
		storage:     store,
	}, nil
}

func (m DefaultBidManager) AddBid(bidder, startingBid, maxBid, incrementAmount string) error {
	start, err := currency.ParseAmount(startingBid)
	if err != nil {
		return errors.Join(errors.New("failed to parse starting bid"), err)
	}

	maxB, err := currency.ParseAmount(maxBid)
	if err != nil {
		return errors.Join(errors.New("failed to parse max bid"), err)
	}

	increment, err := currency.ParseAmount(incrementAmount)
	if err != nil {
		return errors.Join(errors.New("failed to parse increment amount"), err)
	}

	err = m.checkValidBid(start, maxB, increment)
	if err != nil {
		return err
	}

	bid := auction.Bid{
		Bidder:      auction.Bidder(bidder),
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

func (m DefaultBidManager) checkValidBid(startingBid, maxBid, incrementAmount currency.Amount) error {
	if maxBid.Less(startingBid) {
		return &InvalidBidError{message: fmt.Sprintf("starting bid %s cannot be larger than max bid %s", startingBid.String(), maxBid.String())}
	}
	if incrementAmount.Equals(currency.Amount{Dollars: 0, Cents: 0}) {
		return &InvalidBidError{message: fmt.Sprintf("bid increment %s cannot be zero", incrementAmount.String())}
	}
	return nil
}

func (m DefaultBidManager) CalculateWinner() (auction.WinningBid, error) {
	var currentWinner auction.WinningBid

	bids, err := m.storage.GetAllBids()
	if err != nil {
		return auction.WinningBid{}, errors.Join(errors.New("failed to fetch bids"), err)
	}

	if len(bids) == 0 {
		return auction.WinningBid{}, &EmptyBidListError{}
	}

	state := m.initializeCalculation(bids)

	complete := false
	for !complete {
		state = m.calculateBids(bids, state, currentWinner)
		currentWinner = m.currentWinner(bids, state, currentWinner)
		complete = m.isFinished(bids, state, currentWinner)
	}
	return currentWinner, nil
}

func (m DefaultBidManager) initializeCalculation(bids map[auction.Bidder]auction.Bid) bidState {
	state := bidState{}

	for bidder, bid := range bids {
		state[bidder] = bid.StartingBid
	}
	return state
}

func (m DefaultBidManager) calculateBids(bids map[auction.Bidder]auction.Bid, state bidState, currentWinner auction.WinningBid) bidState {
	newState := bidState{}
	for bidder, amount := range state {
		newAmount := amount
		bid := bids[bidder]
		for m.isLessThanCurrentWinner(currentWinner.Amount, newAmount) && m.canStillBid(bid.MaxBid, newAmount, bid.Increment) && bidder != currentWinner.Bidder {
			newAmount = newAmount.Add(bid.Increment)
		}
		newState[bidder] = newAmount
	}
	return newState
}

func (m DefaultBidManager) currentWinner(bids map[auction.Bidder]auction.Bid, state bidState, currentWinner auction.WinningBid) auction.WinningBid {
	highestBidder := currentWinner.Bidder
	for bidder, amount := range state {
		highestBid := state[highestBidder]
		if amount.Greater(highestBid) {
			highestBidder = bidder
		} else if m.isTied(amount, highestBid) {
			highestBidder = m.breakTie(bids[bidder].ID, bids[highestBidder].ID, bidder, highestBidder)
		}
	}
	return auction.WinningBid{
		Bidder: highestBidder,
		Amount: state[highestBidder],
	}
}

func (m DefaultBidManager) isTied(bid, highestBid currency.Amount) bool {
	return bid.Equals(highestBid)
}

func (m DefaultBidManager) breakTie(bidderID, highestBidderID id_generator.EventID, bidder, highestBidder auction.Bidder) auction.Bidder {
	if bidderID < highestBidderID {
		return bidder
	}
	return highestBidder
}

func (m DefaultBidManager) isFinished(bids map[auction.Bidder]auction.Bid, state bidState, currentWinner auction.WinningBid) bool {
	complete := true
	for bidder, amount := range state {
		bid := bids[bidder]
		if m.canStillBid(bid.MaxBid, amount, bid.Increment) && bidder != currentWinner.Bidder {
			complete = false
		}
	}
	return complete
}

func (m DefaultBidManager) isLessThanCurrentWinner(currentWinner, amount currency.Amount) bool {
	return currentWinner.Greater(amount) || currentWinner.Equals(amount)
}

func (m DefaultBidManager) canStillBid(maxBid, amount, increment currency.Amount) bool {
	return maxBid.Sub(amount).Greater(increment) || maxBid.Sub(amount).Equals(increment)
}
