package currency

import "fmt"

type InvalidCurrencyFormatError struct {
	amount string
}

func (e *InvalidCurrencyFormatError) Error() string {
	return fmt.Sprintf("invalid format. Currency must be in the form one of ($0.00, 0.00, $0, 0) . Was given: %s", e.amount)
}
