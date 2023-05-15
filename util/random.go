package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// special function called first of all
func init() {
	rand.Seed(time.Now().UnixNano()) // If Seed is not called, the generator behaves as if seeded by Seed(1).
}

// generates int64 between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates random string which length is n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(5)
}

func RandomBalance() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "JPY", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomDate() time.Time {
	// Generate a random number of seconds within a specific range
	min := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()

	return time.Unix(RandomInt(min, max), 0)
}
