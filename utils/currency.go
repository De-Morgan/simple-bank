package utils

const (
	USD = "USD"
	NGN = "NGN"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, NGN:
		return true
	default:
		return false
	}
}
