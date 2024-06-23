package storage

import (
	"auction/currency"
	"cmp"
	"reflect"
	"slices"
	"testing"
)

type storageTests struct {
	// Make this more generic later
	storeFn func() BidStorer
	t       *testing.T
}

// Add a test for duplicate users
// Add a test for user not found
func (g *storageTests) Run() {
	tests := map[string]func(t *testing.T, store BidStorer){
		"Test Set Get":          testSetGet,
		"Test Set Get Multiple": testSetGetMultiple,
		"Test Get All":          testSetGetAll,
	}
	for name, test := range tests {
		g.t.Run(name, func(t *testing.T) {
			test(t, g.storeFn())
		})
	}
}

func testSetGet(t *testing.T, store BidStorer) {
	expBid := Bid{
		Bidder: Bidder("mockBidder"),
		StartingBid: currency.Amount{
			Dollars: 1,
			Cents:   20,
		},
		MaxBid: currency.Amount{
			Dollars: 5,
			Cents:   6,
		},
		Increment: currency.Amount{
			Dollars: 0,
			Cents:   20,
		},
		ID: 1,
	}
	err := store.SaveBid(expBid)
	if err != nil {
		t.Fatalf("Failed to save bid: %s", err.Error())
	}

	recBid, err := store.GetBid(expBid.Bidder)
	if err != nil {
		t.Fatalf("Failed to get bid: %s", err.Error())
	}

	if !reflect.DeepEqual(expBid, recBid) {
		t.Fatalf("Bids do not match. Expected:\n%#v\nGot:\n%#v", expBid, recBid)
	}
}

func testSetGetMultiple(t *testing.T, store BidStorer) {
	expBids := []Bid{
		{
			Bidder: Bidder("mockBidder"),
			StartingBid: currency.Amount{
				Dollars: 1,
				Cents:   20,
			},
			MaxBid: currency.Amount{
				Dollars: 5,
				Cents:   6,
			},
			Increment: currency.Amount{
				Dollars: 0,
				Cents:   20,
			},
			ID: 1,
		},
		{
			Bidder: Bidder("mockBidder2"),
			StartingBid: currency.Amount{
				Dollars: 3,
				Cents:   45,
			},
			MaxBid: currency.Amount{
				Dollars: 6,
				Cents:   33,
			},
			Increment: currency.Amount{
				Dollars: 1,
				Cents:   5,
			},
			ID: 2,
		},
	}
	for _, bid := range expBids {
		err := store.SaveBid(bid)
		if err != nil {
			t.Fatalf("Failed to save bid: %s", err.Error())
		}
	}

	for _, bid := range expBids {
		recBid, err := store.GetBid(bid.Bidder)
		if err != nil {
			t.Fatalf("Failed to get bid: %s", err.Error())
		}

		if !reflect.DeepEqual(bid, recBid) {
			t.Fatalf("Bids do not match. Expected:\n%#v\nGot:\n%#v", bid, recBid)
		}
	}
}

func testSetGetAll(t *testing.T, store BidStorer) {
	expBids := []Bid{
		{
			Bidder: Bidder("mockBidder"),
			StartingBid: currency.Amount{
				Dollars: 1,
				Cents:   20,
			},
			MaxBid: currency.Amount{
				Dollars: 5,
				Cents:   6,
			},
			Increment: currency.Amount{
				Dollars: 0,
				Cents:   20,
			},
			ID: 1,
		},
		{
			Bidder: Bidder("mockBidder2"),
			StartingBid: currency.Amount{
				Dollars: 3,
				Cents:   45,
			},
			MaxBid: currency.Amount{
				Dollars: 6,
				Cents:   33,
			},
			Increment: currency.Amount{
				Dollars: 1,
				Cents:   5,
			},
			ID: 2,
		},
		{
			Bidder: Bidder("mockBidder3"),
			StartingBid: currency.Amount{
				Dollars: 5,
				Cents:   12,
			},
			MaxBid: currency.Amount{
				Dollars: 8,
				Cents:   45,
			},
			Increment: currency.Amount{
				Dollars: 0,
				Cents:   1,
			},
			ID: 3,
		},
	}
	for _, bid := range expBids {
		err := store.SaveBid(bid)
		if err != nil {
			t.Fatalf("Failed to save bid: %s", err.Error())
		}
	}

	recBids, err := store.GetAllBids()
	if err != nil {
		t.Fatalf("Failed to get bid: %s", err.Error())
	}

	if len(recBids) != len(expBids) {
		t.Fatalf("Expected %d bids, got: %d", len(expBids), len(recBids))
	}

	slices.SortFunc(expBids, func(a, b Bid) int {
		return cmp.Compare(a.ID, b.ID)
	})
	slices.SortFunc(recBids, func(a, b Bid) int {
		return cmp.Compare(a.ID, b.ID)
	})

	for i, recBid := range recBids {
		if !reflect.DeepEqual(recBid, expBids[i]) {
			t.Fatalf("Bids do not match. Expected:\n%#v\nGot:\n%#v", expBids[i], recBid)
		}
	}
}
