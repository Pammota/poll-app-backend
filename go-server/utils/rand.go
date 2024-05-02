package utils

import (
	"math/rand"
	"time"
)

// const upperCaseLetterRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowerCaseLetterRunes = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = lowerCaseLetterRunes[r.Intn(len(lowerCaseLetterRunes))]
	}
	return string(b)
}
