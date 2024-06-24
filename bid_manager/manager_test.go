package bid_manager

import (
	"auction/currency"
	"auction/storage"
	"reflect"
	"testing"
)

// TODO test case for no bids entered
type managerTests struct {
	// Make this more generic later
	managerFn func() (BidManager, error)
	t         *testing.T
}

func (g *managerTests) Run() {
	tests := map[string]func(t *testing.T, manager BidManager){
		"Test General Cases": testGeneralCases,
		"Test Tie Cases":     testTieCases,
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
		winner WinningBid
	}
	testCases := []testCase{
		{
			"Case 1",
			[]bid{
				{"Sasha", "$50.00", "$80.00", "$3.00"},
				{"John", "$60.00", "$82.00", "$2.00"},
				{"Pat", "$55.00", "$85.00", "$5.00"},
			},
			WinningBid{
				storage.Bidder("Pat"),
				currency.Amount{Dollars: 85, Cents: 00},
			},
		},
		{
			"Case 2",
			[]bid{
				{"Riley", "$700.00", "$725.00", "$2.00"},
				{"Morgan", "$599.00", "$725.00", "$15.00"},
				{"Charlie", "$625.00", "$725.00", "$8.00"},
			},
			WinningBid{
				storage.Bidder("Riley"),
				currency.Amount{Dollars: 722, Cents: 00},
			},
		},
		{
			"Case 3",
			[]bid{
				{"Alex", "$2500.00", "$3000.00", "$500.00"},
				{"Jesse", "$2800.00", "$3100.00", "$201.00"},
				{"Drew", "$2501.00", "$3200.00", "$247.00"},
			},
			WinningBid{
				storage.Bidder("Jesse"),
				currency.Amount{Dollars: 3001, Cents: 00},
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
		winner WinningBid
	}
	testCases := []testCase{
		{
			"Case 1",
			[]bid{
				{"Sasha", "$50.00", "$80.00", "$3.00"},
				{"John", "$50.00", "$80.00", "$3.00"},
				{"Pat", "$50.00", "$80.00", "$3.00"},
			},
			WinningBid{
				storage.Bidder("Sasha"),
				currency.Amount{Dollars: 80, Cents: 00},
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
