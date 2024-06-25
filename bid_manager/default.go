package bid_manager

import (
	"auction/auction"
	"auction/currency"
	"auction/id_generator"
	"auction/storage"
	"errors"
	"fmt"
)

// bidState stores each bidders current bid value as the winner is determined
type bidState map[auction.Bidder]currency.Amount

// defaultBidManager implements the BidManager interface and can be provided with different implementations for storage
// and ID generation
type defaultBidManager struct {
	idGenerator id_generator.IDGenerator
	storage     storage.BidStorer
}

func NewDefaultBidManager(idGenerator id_generator.IDGenerator, store storage.BidStorer) (BidManager, error) {
	return &defaultBidManager{
		idGenerator: idGenerator,
		storage:     store,
	}, nil
}

// AddBid takes a bid entry as strings, then parses and saves them to be used later to calculate the winning bid.
func (m defaultBidManager) AddBid(bidder, startingBid, maxBid, incrementAmount string) error {
	start, err := currency.ParseAmount(startingBid)
	if err != nil {
		return errors.Join(&InvalidBidError{message: "failed to parse starting bid"}, err)
	}

	maxB, err := currency.ParseAmount(maxBid)
	if err != nil {
		return errors.Join(&InvalidBidError{message: "failed to parse max bid"}, err)
	}

	increment, err := currency.ParseAmount(incrementAmount)
	if err != nil {
		return errors.Join(&InvalidBidError{message: "failed to parse increment amount"}, err)
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

// checkValidBid ensures that valid, non-zero or negative values, are given for the bid
func (m defaultBidManager) checkValidBid(startingBid, maxBid, incrementAmount currency.Amount) error {
	if maxBid.Less(startingBid) {
		return &InvalidBidError{message: fmt.Sprintf("starting bid %s cannot be larger than max bid %s", startingBid.String(), maxBid.String())}
	}
	if incrementAmount.Less(currency.Amount{Dollars: 0, Cents: 1}) {
		return &InvalidBidError{message: fmt.Sprintf("bid increment %s cannot be less than 1 cent", incrementAmount.String())}
	}
	if startingBid.Less(currency.Amount{Dollars: 0, Cents: 1}) {
		return &InvalidBidError{message: fmt.Sprintf("starting bid %s cannot be less than 1 cent", incrementAmount.String())}
	}
	if maxBid.Less(currency.Amount{Dollars: 0, Cents: 1}) {
		return &InvalidBidError{message: fmt.Sprintf("max increment %s cannot be less than 1 cent", incrementAmount.String())}
	}
	return nil
}

// CalculateWinner iterates through all of the provided bids to determine what the winning bid will be. It does so in
// 'rounds', where each round the non-current winners have their bid incremented until it is greater than the winning
// bid or until they can no longer bid without exceeding their max bid. Once no more bids can be incremented to beat the
// current winner, it returns the WinningBid which contains the winners name and bid amount.
func (m defaultBidManager) CalculateWinner() (auction.WinningBid, error) {
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

// initializeCalculation sets the initial rounds of bids to each person starting bid amount.
func (m defaultBidManager) initializeCalculation(bids map[auction.Bidder]auction.Bid) bidState {
	state := bidState{}

	for bidder, bid := range bids {
		state[bidder] = bid.StartingBid
	}
	return state
}

// calculateBids checks to see if each person is bidding under the current winner and is still able to bid. It will
// then increment their current amount until it exceeds the winner but is still under their max bid amount.
func (m defaultBidManager) calculateBids(bids map[auction.Bidder]auction.Bid, state bidState, currentWinner auction.WinningBid) bidState {
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

// currentWinner checks to see which bidder is the current winner to be used for the next round of bids or as the final
// winner
func (m defaultBidManager) currentWinner(bids map[auction.Bidder]auction.Bid, state bidState, currentWinner auction.WinningBid) auction.WinningBid {
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

// isTied checks to see if the persons bid is tied with the current highest bid
func (m defaultBidManager) isTied(bid, highestBid currency.Amount) bool {
	return bid.Equals(highestBid)
}

// breakTie breaks a tie based on who has the lowest ID, which signifies that they entered their bid first
func (m defaultBidManager) breakTie(bidderID, highestBidderID id_generator.EventID, bidder, highestBidder auction.Bidder) auction.Bidder {
	if bidderID < highestBidderID {
		return bidder
	}
	return highestBidder
}

// isFinished checks to see if there are any bids that can still be placed without exceeding the persons max bid
func (m defaultBidManager) isFinished(bids map[auction.Bidder]auction.Bid, state bidState, currentWinner auction.WinningBid) bool {
	complete := true
	for bidder, amount := range state {
		bid := bids[bidder]
		if m.canStillBid(bid.MaxBid, amount, bid.Increment) && bidder != currentWinner.Bidder {
			complete = false
		}
	}
	return complete
}

// isLessThanCurrentWinner checks to see if the current bidder is bidding less than the winning bid
func (m defaultBidManager) isLessThanCurrentWinner(currentWinner, amount currency.Amount) bool {
	return currentWinner.Greater(amount) || currentWinner.Equals(amount)
}

// canStillBid checks to see if the different between the person max bid is less than or equal to
// the max amount the person wants to bid, signifying that their bid can still be increased
func (m defaultBidManager) canStillBid(maxBid, amount, increment currency.Amount) bool {
	return maxBid.Sub(amount).Greater(increment) || maxBid.Sub(amount).Equals(increment)
}
