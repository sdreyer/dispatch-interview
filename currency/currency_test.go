package currency

import (
	"reflect"
	"testing"
)

func TestAmount_Add(t *testing.T) {
	type testCase struct {
		name     string
		a1       Amount
		a2       Amount
		expected Amount
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}},
		{"OneZero", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 1, Cents: 55}},
		{"Simple", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 5, Cents: 29}, Amount{Dollars: 6, Cents: 84}},
		{"CentRollover", Amount{Dollars: 1, Cents: 99}, Amount{Dollars: 5, Cents: 2}, Amount{Dollars: 7, Cents: 1}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			total := test.a1.Add(test.a2)
			if !reflect.DeepEqual(total, test.expected) {
				t.Errorf("Expected: %#v, got: %#v", test.expected, total)
			}
		})
	}
}

func TestAmount_Sub(t *testing.T) {
	type testCase struct {
		name     string
		a1       Amount
		a2       Amount
		expected Amount
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}},
		{"SubZero", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 1, Cents: 55}},
		{"Simple", Amount{Dollars: 1, Cents: 00}, Amount{Dollars: 0, Cents: 29}, Amount{Dollars: 0, Cents: 71}},
		{"CentRollover", Amount{Dollars: 1, Cents: 99}, Amount{Dollars: 5, Cents: 2}, Amount{Dollars: -3, Cents: -3}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			total := test.a1.Sub(test.a2)
			if !reflect.DeepEqual(total, test.expected) {
				t.Errorf("Expected: %#v, got: %#v", test.expected, total)
			}
		})
	}
}

func TestAmount_Equals(t *testing.T) {
	type testCase struct {
		name     string
		a1       Amount
		a2       Amount
		expected bool
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}, true},
		{"First Larger", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 0, Cents: 0}, false},
		{"Second Larger", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 5, Cents: 29}, false},
		{"Simple", Amount{Dollars: 5, Cents: 99}, Amount{Dollars: 5, Cents: 99}, true},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			isEqual := test.a1.Equals(test.a2)
			if !reflect.DeepEqual(isEqual, test.expected) {
				t.Errorf("Expected equal to be %v, but was: %v", test.expected, isEqual)
			}
		})
	}
}

func TestAmount_ToString(t *testing.T) {
	type testCase struct {
		name     string
		a1       Amount
		expected string
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, "$0.00"},
		{"One Cent", Amount{Dollars: 0, Cents: 1}, "$0.01"},
		{"Multiple Cents", Amount{Dollars: 0, Cents: 15}, "$0.15"},
		{"Dollar No Cents", Amount{Dollars: 1, Cents: 0}, "$1.00"},
		{"Dollar And Cents", Amount{Dollars: 1, Cents: 33}, "$1.33"},
		{"Negative", Amount{Dollars: -1, Cents: -33}, "-$1.33"},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			stringVal := test.a1.String()
			if test.expected != stringVal {
				t.Errorf("Expected value to be %s, but was: %s", test.expected, stringVal)
			}
		})
	}
}

func TestAmount_ParseString(t *testing.T) {
	type testCase struct {
		name        string
		given       string
		expected    Amount
		shouldError bool
	}
	testCases := []testCase{
		{"Zero", "$0.00", Amount{Dollars: 0, Cents: 0}, false},
		{"Zero No Dollar Sign", "0.00", Amount{Dollars: 0, Cents: 0}, false},
		{"Only Cents", "$0.55", Amount{Dollars: 0, Cents: 55}, false},
		{"Only Dollar", "$1.00", Amount{Dollars: 1, Cents: 0}, false},
		{"No Decimal Cents", "$1", Amount{Dollars: 1, Cents: 0}, false},
		{"No Dollar", "$.55", Amount{Dollars: 0, Cents: 55}, true},
		{"Dollar and Cents", "$5.99", Amount{Dollars: 5, Cents: 99}, false},
		{"Partial Cents", "$0.5", Amount{Dollars: 0, Cents: 50}, false},
		{"Empty String", "", Amount{}, true},
		{"Extra Digits", "$0.55555", Amount{}, true},
		{"Negative", "-$1.55", Amount{Dollars: -1, Cents: -55}, false},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			amount, err := ParseAmount(test.given)
			if err != nil && !test.shouldError {
				t.Fatalf("operation failed with: %s", err.Error())
			} else if err == nil && test.shouldError {
				t.Fatalf("operation should have failed but didn't")
			} else if err != nil && test.shouldError {
				return
			}

			if test.expected != amount {
				t.Fatalf("Expected value to be %s, but was: %s", test.expected, amount)
			}
		})
	}
}

func TestAmount_Less(t *testing.T) {
	type testCase struct {
		name     string
		a1       Amount
		a2       Amount
		expected bool
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}, false},
		{"First Larger", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 0, Cents: 0}, false},
		{"Second Larger", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 5, Cents: 29}, true},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			isLess := test.a1.Less(test.a2)
			if !reflect.DeepEqual(isLess, test.expected) {
				t.Errorf("Expected less to be %v, but was: %v", test.expected, isLess)
			}
		})
	}
}

func TestAmount_Greater(t *testing.T) {
	type testCase struct {
		name     string
		a1       Amount
		a2       Amount
		expected bool
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}, false},
		{"First Larger", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 0, Cents: 0}, true},
		{"Second Larger", Amount{Dollars: 1, Cents: 55}, Amount{Dollars: 5, Cents: 29}, false},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			isGreater := test.a1.Greater(test.a2)
			if !reflect.DeepEqual(isGreater, test.expected) {
				t.Errorf("Expected greater to be %v, but was: %v", test.expected, isGreater)
			}
		})
	}
}

func TestAmount_abs(t *testing.T) {
	type testCase struct {
		name     string
		a        int64
		expected int64
	}
	testCases := []testCase{
		{"Zero", 0, 0},
		{"Positive", 1, 1},
		{"Negative", -1, 1},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			absVal := abs(test.a)
			if absVal != test.expected {
				t.Fatalf("Expected value %d, got %d", test.expected, absVal)
			}
		})
	}
}

func TestAmount_CurrentAbs(t *testing.T) {
	type testCase struct {
		name     string
		a        Amount
		expected Amount
	}
	testCases := []testCase{
		{"Zero", Amount{Dollars: 0, Cents: 0}, Amount{Dollars: 0, Cents: 0}},
		{"Positive", Amount{Dollars: 1, Cents: 1}, Amount{Dollars: 1, Cents: 1}},
		{"Negative", Amount{Dollars: -1, Cents: -1}, Amount{Dollars: 1, Cents: 1}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			absVal := test.a.Abs()
			if !reflect.DeepEqual(absVal, test.expected) {
				t.Fatalf("Expected value %v, got %v", test.expected, absVal)
			}
		})
	}
}
