package bid_manager

import (
	"auction/auction"
	"auction/currency"
	"reflect"
	"testing"
)

type managerTests struct {
	managerFn func() (BidManager, error)
	t         *testing.T
}

func (g *managerTests) Run() {
	tests := map[string]func(t *testing.T, manager BidManager){
		"Test General Cases":                 testGeneralCases,
		"Test Tie Cases":                     testTieCases,
		"Test Empty Bid List":                testEmptyBidList,
		"Test Starting Bid Greater Than Max": testStartingBidGreaterThanMax,
		"Test Zero Increment":                testZeroIncrement,
	}
	for name, test := range tests {
		g.t.Run(name, func(t *testing.T) {
			manager, err := g.managerFn()
			if err != nil {
				t.Fatalf("could not initialize manager: %s", err.Error())
			}
			test(t, manager)
		})
	}
}

/*
Provided Test Case 1:
         Initial Bid         Max Bid         Bid Increment
Sasha      $50.00            $80.00             $3.00
John       $60.00            $82.00             $2.00
Pat        $55.00            $85.00             $5.00

          Sasha      John       Pat       Current Winner
Round 1   $50.00     $60.00     $55.00    John
Round 2   $62.00     $60.00     $65.00    Pat
Round 3   $68.00     $66.00     $65.00    Sasha
Round 4   $68.00     $70.00     $70.00    John + Pat
Round 5   $71.00     $72.00     $75.00    Pat
Round 6   $77.00     $76.00     $75.00    Sasha
Round 7   $77.00     $78.00     $80.00    Pat
Round 8   $80.00     $82.00     $80.00    John
Round 9   $80.00     $82.00     $85.00    Pat
Winner is Pat @ $85.00

Provided Test Case 2:
         Initial Bid         Max Bid         Bid Increment
Riley      $700.00           $725.00            $2.00
Morgan     $599.00           $725.00            $15.00
Charlie    $625.00           $725.00            $8.00

          Riley      Morgan     Charlie     Current Winner
Round 1   $700.00    $599.00    $625.00     Riley
Round 2   $700.00    $704.00    $705.00     Charlie
Round 3   $706.00    $719.00    $705.00     Morgan
Round 4   $720.00    $719.00    $721.00     Charlie
Round 5   $722.00    $719.00    $721.00     Riley
Winner is Riley @ 722.00

Provided Test Case 3:
         Initial Bid         Max Bid         Bid Increment
Alex      $2500.00           $3000.00          $500.00
Jesse     $2800.00           $3100.00          $201.00
Drew      $2501.00           $3200.00          $247.00

          Alex       Jesse      Drew        Current Winner
Round 1   $2500.00   $2800.00   $2501.00    Jesse
Round 2   $3000.00   $2800.00   $2995.00    Alex
Round 3   $3000.00   $3001.00   $2995.00    Jesse
Winner is Jesse @ 3001.00
*/

func testGeneralCases(t *testing.T, manager BidManager) {
	type bid struct {
		bidder     string
		initialBid string
		maxBid     string
		increment  string
	}
	type testCase struct {
		name   string
		bids   []bid
		winner auction.WinningBid
	}
	testCases := []testCase{
		{
			name: "Case 1",
			bids: []bid{
				{"Sasha", "$50.00", "$80.00", "$3.00"},
				{"John", "$60.00", "$82.00", "$2.00"},
				{"Pat", "$55.00", "$85.00", "$5.00"},
			},
			winner: auction.WinningBid{
				Bidder: auction.Bidder("Pat"),
				Amount: currency.Amount{Dollars: 85, Cents: 00},
			},
		},
		{
			name: "Case 2",
			bids: []bid{
				{"Riley", "$700.00", "$725.00", "$2.00"},
				{"Morgan", "$599.00", "$725.00", "$15.00"},
				{"Charlie", "$625.00", "$725.00", "$8.00"},
			},
			winner: auction.WinningBid{
				Bidder: auction.Bidder("Riley"),
				Amount: currency.Amount{Dollars: 722, Cents: 00},
			},
		},
		{
			name: "Case 3",
			bids: []bid{
				{"Alex", "$2500.00", "$3000.00", "$500.00"},
				{"Jesse", "$2800.00", "$3100.00", "$201.00"},
				{"Drew", "$2501.00", "$3200.00", "$247.00"},
			},
			winner: auction.WinningBid{
				Bidder: auction.Bidder("Jesse"),
				Amount: currency.Amount{Dollars: 3001, Cents: 00},
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			for _, bid := range test.bids {
				err := manager.AddBid(bid.bidder, bid.initialBid, bid.maxBid, bid.increment)
				if err != nil {
					t.Fatalf("Failed to add bid: %s", err.Error())
				}
			}
			recWinner, err := manager.CalculateWinner()
			if err != nil {
				t.Fatalf("Failed to calculate winner: %s", err.Error())
			}
			if !reflect.DeepEqual(recWinner, test.winner) {
				t.Fatalf("Expected %#v, got %#v", test.winner, recWinner)
			}
		})
	}
}

func testTieCases(t *testing.T, manager BidManager) {
	type bid struct {
		bidder     string
		initialBid string
		maxBid     string
		increment  string
	}
	type testCase struct {
		name   string
		bids   []bid
		winner auction.WinningBid
	}
	testCases := []testCase{
		{
			name: "Case 1",
			bids: []bid{
				{"Sasha", "$50.00", "$80.00", "$3.00"},
				{"John", "$50.00", "$80.00", "$3.00"},
				{"Pat", "$50.00", "$80.00", "$3.00"},
			},
			winner: auction.WinningBid{
				Bidder: auction.Bidder("Sasha"),
				Amount: currency.Amount{Dollars: 80, Cents: 00},
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			for _, bid := range test.bids {
				err := manager.AddBid(bid.bidder, bid.initialBid, bid.maxBid, bid.increment)
				if err != nil {
					t.Fatalf("Failed to add bid: %s", err.Error())
				}
			}
			recWinner, err := manager.CalculateWinner()
			if err != nil {
				t.Fatalf("Failed to calculate winner: %s", err.Error())
			}
			if !reflect.DeepEqual(recWinner, test.winner) {
				t.Fatalf("Expected %#v, got %#v", test.winner, recWinner)
			}
		})
	}
}

func testEmptyBidList(t *testing.T, manager BidManager) {
	_, err := manager.CalculateWinner()
	if err == nil {
		t.Fatalf("Expected EmptyBidListError and did not receive one")
	} else if _, ok := err.(*EmptyBidListError); !ok {
		t.Fatalf("Expected EmptyBidListError but got: %#v", err)
	}
}

func testStartingBidGreaterThanMax(t *testing.T, manager BidManager) {
	err := manager.AddBid("mockBidder", "$100", "$5", "$0.20")
	if err == nil {
		t.Fatalf("Expected InvalidBidError and did not receive one")
	} else if _, ok := err.(*InvalidBidError); !ok {
		t.Fatalf("Expected InvalidBidError but got %#v", err)
	}

}

func testZeroIncrement(t *testing.T, manager BidManager) {
	err := manager.AddBid("mockBidder", "$5", "$20", "$0")
	if err == nil {
		t.Fatalf("Expected InvalidBidError and did not receive one")
	} else if _, ok := err.(*InvalidBidError); !ok {
		t.Fatalf("Expected InvalidBidError but got %#v", err)
	}

}
