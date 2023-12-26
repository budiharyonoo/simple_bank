package utils

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt Ggenerate random int from min - max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabets = "abcdefghijklmnopqrstuvwxyz"

// RandomString generate random string based on length
func RandomString(stringLength int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < stringLength; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generate random owner name for test
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generate random amount of money for test
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generate random currency code for test
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "IDR", "SGD"}
	return currencies[rand.Intn(len(currencies))]
}
