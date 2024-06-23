package currency

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Amount struct {
	Dollars uint64
	Cents   uint64
}

func (a Amount) String() string {
	return fmt.Sprintf("$%d.%02d", a.Dollars, a.Cents)
}

func (a Amount) Add(amt Amount) Amount {
	totalCents := a.Cents + amt.Cents
	totalDollars := a.Dollars + amt.Dollars + totalCents/100
	return Amount{
		Dollars: totalDollars,
		Cents:   totalCents % 100,
	}
}

func (a Amount) Sub(amt Amount) (Amount, error) {
	totalCentsA := a.Dollars*100 + a.Cents
	totalCentsB := amt.Dollars*100 + amt.Cents

	if totalCentsA < totalCentsB {
		return Amount{}, fmt.Errorf("failed to subtract %s from %s. Operation results in a negative amount", amt.String(), a.String())
	}

	diffCents := totalCentsA - totalCentsB
	return Amount{
		Dollars: diffCents / 100,
		Cents:   diffCents % 100,
	}, nil
}

func (a Amount) Equals(amt Amount) bool {
	return a.Dollars == amt.Dollars && a.Cents == amt.Cents
}

func ParseAmount(s string) (Amount, error) {
	s = strings.TrimPrefix(s, "$")

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return Amount{}, errors.New("invalid amount format")
	}

	dollars, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return Amount{}, fmt.Errorf("invalid dollars value: %v", err)
	}

	var cents uint64

	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return Amount{}, errors.New("cents part should not have more than two digits")
		}
		centsString := parts[1] + strings.Repeat("0", 2-len(parts[1]))
		cents, err = strconv.ParseUint(centsString, 10, 64)
		if err != nil {
			return Amount{}, fmt.Errorf("invalid cents value: %v", err)
		}
	}

	return Amount{Dollars: dollars, Cents: cents}, nil
}
