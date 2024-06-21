package utils

import (
	"math/rand"
	"strings"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) (result int64) {
	result = min + random.Int63n(max-min+1)
	return
}

func RandomString(length int) (result string) {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < length; i++ {
		c := alphabet[random.Intn(k)]
		sb.WriteByte(c)
	}
	result = sb.String()
	return
}

func RandomCurrency() string {
	c := []string{"USD", "EUR", "NGN"}
	return c[random.Intn(len(c))]
}

func RandomName() string {
	return RandomString(8)
}
func RandomMoney() int64 {
	return RandomInt(1000, 10000)
}
