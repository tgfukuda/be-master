package util

// supported currency
const (
	USD = "USD"
	EUR = "EUR"
	JPY = "JPY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, JPY:
		return true
	}
	return false
}
