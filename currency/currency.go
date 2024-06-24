package currency

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Amount struct {
	Dollars int64
	Cents   int64
}

func (a Amount) String() string {
	sign := ""
	if a.Dollars < 0 || a.Cents < 0 {
		sign = "-"
		a = a.Abs()
	}
	return fmt.Sprintf("%s$%d.%02d", sign, a.Dollars, a.Cents)
}

func (a Amount) Abs() Amount {
	return Amount{
		Dollars: abs(a.Dollars),
		Cents:   abs(a.Cents),
	}
}

func abs(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

func (a Amount) Add(amt Amount) Amount {
	totalCents := a.Cents + amt.Cents
	totalDollars := a.Dollars + amt.Dollars
	// Adjust for cent overflow
	if totalCents >= 100 {
		totalDollars++
		totalCents -= 100
	} else if totalCents < 0 {
		totalDollars--
		totalCents += 100
	}
	return Amount{
		Dollars: totalDollars,
		Cents:   totalCents,
	}
}

func (a Amount) Sub(amt Amount) Amount {
	totalCentsA := a.Dollars*100 + a.Cents
	totalCentsB := amt.Dollars*100 + amt.Cents

	diffCents := totalCentsA - totalCentsB
	return Amount{
		Dollars: diffCents / 100,
		Cents:   diffCents % 100,
	}
}

func (a Amount) Equals(amt Amount) bool {
	return a.Dollars == amt.Dollars && a.Cents == amt.Cents
}

func (a Amount) Less(amt Amount) bool {
	if a.Dollars < amt.Dollars {
		return true
	} else if a.Dollars == amt.Dollars && a.Cents < amt.Cents {
		return true
	}
	return false
}

func (a Amount) Greater(amt Amount) bool {
	if a.Dollars > amt.Dollars {
		return true
	} else if a.Dollars == amt.Dollars && a.Cents > amt.Cents {
		return true
	}
	return false
}

func ParseAmount(s string) (Amount, error) {
	s = strings.TrimPrefix(s, "$")
	sign := int64(1)
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = strings.TrimPrefix(s, "-")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return Amount{}, errors.New("invalid amount format")
	}

	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Amount{}, fmt.Errorf("invalid dollars value: %v", err)
	}

	var cents int64

	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return Amount{}, errors.New("cents part should not have more than two digits")
		}
		centsString := parts[1] + strings.Repeat("0", 2-len(parts[1]))
		cents, err = strconv.ParseInt(centsString, 10, 64)
		if err != nil {
			return Amount{}, fmt.Errorf("invalid cents value: %v", err)
		}
	}

	return Amount{
		Dollars: sign * dollars,
		Cents:   sign * cents,
	}, nil
}
