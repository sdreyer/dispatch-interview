package bid_manager

import (
	"auction/auction"
	"auction/currency"
	"auction/id_generator"
	"auction/storage"
	"reflect"
	"testing"
)

func WithDefaultBidManager() func() (BidManager, error) {
	return func() (BidManager, error) {
		return NewDefaultBidManager(id_generator.NewMemoryIDGenerator(), storage.NewMemoryBidStorage())
	}
}

func TestManager(t *testing.T) {

	tests := managerTests{
		managerFn: WithDefaultBidManager(),
		t:         t,
	}
	tests.Run()
}

func TestInitializeCalculation(t *testing.T) {
	manager := &DefaultBidManager{
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	bids := map[auction.Bidder]auction.Bid{
		auction.Bidder("bidder1"): {
			auction.Bidder("bidder1"),
			currency.Amount{Dollars: 1, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(1),
		},
		auction.Bidder("bidder2"): {
			auction.Bidder("bidder2"),
			currency.Amount{Dollars: 2, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(2),
		},
		auction.Bidder("bidder3"): {
			auction.Bidder("bidder3"),
			currency.Amount{Dollars: 3, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}
	expectedState := bidState{
		auction.Bidder("bidder1"): currency.Amount{Dollars: 1, Cents: 20},
		auction.Bidder("bidder2"): currency.Amount{Dollars: 2, Cents: 20},
		auction.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)
	if !reflect.DeepEqual(expectedState, state) {
		t.Fatalf("Expected state to be %#v, got %#v", expectedState, state)
	}
}

func TestCalculateBids(t *testing.T) {
	manager := &DefaultBidManager{
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	bids := map[auction.Bidder]auction.Bid{
		auction.Bidder("bidder1"): {
			auction.Bidder("bidder1"),
			currency.Amount{Dollars: 1, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 75},
			id_generator.EventID(1),
		},
		auction.Bidder("bidder2"): {
			auction.Bidder("bidder2"),
			currency.Amount{Dollars: 2, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 30},
			id_generator.EventID(2),
		},
		auction.Bidder("bidder3"): {
			auction.Bidder("bidder3"),
			currency.Amount{Dollars: 3, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}
	expectedState := bidState{
		auction.Bidder("bidder1"): currency.Amount{Dollars: 3, Cents: 45},
		auction.Bidder("bidder2"): currency.Amount{Dollars: 3, Cents: 40},
		auction.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
	}

	currentWinner := auction.WinningBid{
		Bidder: auction.Bidder("bidder3"),
		Amount: currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)

	state = manager.calculateBids(bids, state, currentWinner)
	if !reflect.DeepEqual(expectedState, state) {
		t.Fatalf("Expected state to be \n%#v\ngot \n%#v", expectedState, state)
	}
}

func TestCurrentWinner(t *testing.T) {
	manager := &DefaultBidManager{
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	bids := map[auction.Bidder]auction.Bid{
		auction.Bidder("bidder1"): {
			auction.Bidder("bidder1"),
			currency.Amount{Dollars: 1, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 75},
			id_generator.EventID(1),
		},
		auction.Bidder("bidder2"): {
			auction.Bidder("bidder2"),
			currency.Amount{Dollars: 2, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 30},
			id_generator.EventID(2),
		},
		auction.Bidder("bidder3"): {
			auction.Bidder("bidder3"),
			currency.Amount{Dollars: 3, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}
	expectedState := bidState{
		auction.Bidder("bidder1"): currency.Amount{Dollars: 3, Cents: 45},
		auction.Bidder("bidder2"): currency.Amount{Dollars: 3, Cents: 40},
		auction.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
	}

	expectedWinner := auction.WinningBid{
		Bidder: auction.Bidder("bidder1"), Amount: currency.Amount{Dollars: 3, Cents: 45},
	}

	currentWinner := auction.WinningBid{
		Bidder: auction.Bidder("bidder3"),
		Amount: currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)

	state = manager.calculateBids(bids, state, currentWinner)

	currentWinner = manager.currentWinner(bids, state, currentWinner)

	if !reflect.DeepEqual(expectedWinner, currentWinner) {
		t.Fatalf("Expected current winner to be %#v, got %#v", expectedState, state)
	}
}

func TestIfFinished(t *testing.T) {
	manager := &DefaultBidManager{
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	bids := map[auction.Bidder]auction.Bid{
		auction.Bidder("bidder1"): {
			auction.Bidder("bidder1"),
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 75},
			id_generator.EventID(1),
		},
		auction.Bidder("bidder2"): {
			auction.Bidder("bidder2"),
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 30},
			id_generator.EventID(2),
		},
		auction.Bidder("bidder3"): {
			auction.Bidder("bidder3"),
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}

	currentWinner := auction.WinningBid{
		Bidder: auction.Bidder("bidder3"),
		Amount: currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)

	state = manager.calculateBids(bids, state, currentWinner)

	currentWinner = manager.currentWinner(bids, state, currentWinner)

	finished := manager.isFinished(bids, state, currentWinner)

	if !finished {
		t.Fatalf("Expected calculation to be complete and was not")
	}
}
