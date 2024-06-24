package bid_manager

import (
	"auction/currency"
	"auction/id_generator"
	"auction/storage"
	"reflect"
	"testing"
)

func TestManager(t *testing.T) {

	tests := managerTests{
		managerFn: NewMemoryBidManager,
		t:         t,
	}
	tests.Run()
}

func TestInitializeCalculation(t *testing.T) {
	manager := &MemoryBidManager{
		// These should probably have an option to throw an error on init
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	// TODO check that starting bid is not greater than max bid
	bids := map[storage.Bidder]storage.Bid{
		storage.Bidder("bidder1"): {
			storage.Bidder("bidder1"),
			currency.Amount{Dollars: 1, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(1),
		},
		storage.Bidder("bidder2"): {
			storage.Bidder("bidder2"),
			currency.Amount{Dollars: 2, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(2),
		},
		storage.Bidder("bidder3"): {
			storage.Bidder("bidder3"),
			currency.Amount{Dollars: 3, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}
	expectedState := bidState{
		storage.Bidder("bidder1"): currency.Amount{Dollars: 1, Cents: 20},
		storage.Bidder("bidder2"): currency.Amount{Dollars: 2, Cents: 20},
		storage.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)
	if !reflect.DeepEqual(expectedState, state) {
		t.Fatalf("Expected state to be %#v, got %#v", expectedState, state)
	}
}

func TestCalculateBids(t *testing.T) {
	manager := &MemoryBidManager{
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	// TODO check that starting bid is not greater than max bid
	// TODO check that increment is not zero
	bids := map[storage.Bidder]storage.Bid{
		storage.Bidder("bidder1"): {
			storage.Bidder("bidder1"),
			currency.Amount{Dollars: 1, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 75},
			id_generator.EventID(1),
		},
		storage.Bidder("bidder2"): {
			storage.Bidder("bidder2"),
			currency.Amount{Dollars: 2, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 30},
			id_generator.EventID(2),
		},
		storage.Bidder("bidder3"): {
			storage.Bidder("bidder3"),
			currency.Amount{Dollars: 3, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}
	expectedState := bidState{
		storage.Bidder("bidder1"): currency.Amount{Dollars: 3, Cents: 45},
		storage.Bidder("bidder2"): currency.Amount{Dollars: 3, Cents: 40},
		storage.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
	}

	currentWinner := WinningBid{
		storage.Bidder("bidder3"),
		currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)

	state = manager.calculateBids(bids, state, currentWinner)
	if !reflect.DeepEqual(expectedState, state) {
		t.Fatalf("Expected state to be %#v, got %#v", expectedState, state)
	}
}

func TestCurrentWinner(t *testing.T) {
	manager := &MemoryBidManager{
		idGenerator: id_generator.NewMemoryIDGenerator(),
		storage:     storage.NewMemoryBidStorage(),
	}

	// TODO check that starting bid is not greater than max bid
	// TODO check that increment is not zero
	bids := map[storage.Bidder]storage.Bid{
		storage.Bidder("bidder1"): {
			storage.Bidder("bidder1"),
			currency.Amount{Dollars: 1, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 75},
			id_generator.EventID(1),
		},
		storage.Bidder("bidder2"): {
			storage.Bidder("bidder2"),
			currency.Amount{Dollars: 2, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 0, Cents: 30},
			id_generator.EventID(2),
		},
		storage.Bidder("bidder3"): {
			storage.Bidder("bidder3"),
			currency.Amount{Dollars: 3, Cents: 20},
			currency.Amount{Dollars: 5, Cents: 20},
			currency.Amount{Dollars: 1, Cents: 00},
			id_generator.EventID(3),
		},
	}
	expectedState := bidState{
		storage.Bidder("bidder1"): currency.Amount{Dollars: 3, Cents: 45},
		storage.Bidder("bidder2"): currency.Amount{Dollars: 3, Cents: 40},
		storage.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
	}

	expectedWinner := WinningBid{
		storage.Bidder("bidder1"), currency.Amount{Dollars: 3, Cents: 45},
	}

	currentWinner := WinningBid{
		storage.Bidder("bidder3"),
		currency.Amount{Dollars: 3, Cents: 20},
	}

	state := manager.initializeCalculation(bids)

	state = manager.calculateBids(bids, state, currentWinner)

	currentWinner = manager.currentWinner(bids, state, currentWinner)

	if !reflect.DeepEqual(expectedWinner, currentWinner) {
		t.Fatalf("Expected current winner to be %#v, got %#v", expectedState, state)
	}
}

//func TestIfFinished(t *testing.T) {
//	manager := &MemoryBidManager{
//		idGenerator: id_generator.NewMemoryIDGenerator(),
//		storage:     storage.NewMemoryBidStorage(),
//	}
//
//	// TODO check that starting bid is not greater than max bid
//	// TODO check that increment is not zero
//	bids := map[storage.Bidder]storage.Bid{
//		storage.Bidder("bidder1"): {
//			storage.Bidder("bidder1"),
//			currency.Amount{Dollars: 1, Cents: 20},
//			currency.Amount{Dollars: 5, Cents: 20},
//			currency.Amount{Dollars: 0, Cents: 75},
//			id_generator.EventID(1),
//		},
//		storage.Bidder("bidder2"): {
//			storage.Bidder("bidder2"),
//			currency.Amount{Dollars: 2, Cents: 20},
//			currency.Amount{Dollars: 5, Cents: 20},
//			currency.Amount{Dollars: 0, Cents: 30},
//			id_generator.EventID(2),
//		},
//		storage.Bidder("bidder3"): {
//			storage.Bidder("bidder3"),
//			currency.Amount{Dollars: 3, Cents: 20},
//			currency.Amount{Dollars: 5, Cents: 20},
//			currency.Amount{Dollars: 1, Cents: 00},
//			id_generator.EventID(3),
//		},
//	}
//	expectedState := bidState{
//		storage.Bidder("bidder1"): currency.Amount{Dollars: 3, Cents: 45},
//		storage.Bidder("bidder2"): currency.Amount{Dollars: 3, Cents: 40},
//		storage.Bidder("bidder3"): currency.Amount{Dollars: 3, Cents: 20},
//	}
//
//	expectedWinner := WinningBid{
//		storage.Bidder("bidder1"), currency.Amount{Dollars: 3, Cents: 45},
//	}
//
//	currentWinner := WinningBid{
//		storage.Bidder("bidder3"),
//		currency.Amount{Dollars: 3, Cents: 20},
//	}
//
//	state := manager.initializeCalculation(bids)
//
//	state = manager.calculateBids(bids, state, currentWinner)
//
//	currentWinner = manager.currentWinner(state, currentWinner)
//
//	finished := manager.isFinished(bids, state, currentWinner)
//
//	if !reflect.DeepEqual(expectedWinner, currentWinner) {
//		t.Fatalf("Expected current winner to be %#v, got %#v", expectedState, state)
//	}
//}
