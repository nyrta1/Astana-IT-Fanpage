package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	verificationCodeLength = 6
)

func GenerateVerificationCode() string {
	code := make([]byte, verificationCodeLength)

	characters := "0123456789"

	for i := 0; i < verificationCodeLength; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		code[i] = characters[randomIndex.Int64()]
	}

	return string(code)
}
