package utils

const (
	USD = "USD"
	EUR = "EUR"
	SGD = "SGD"
	IDR = "IDR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, SGD, IDR:
		return true
	}
	return false
}
