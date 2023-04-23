package main

import (
	"crypto/rand"
	"math/big"
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	for i := 0; i < length; i++ {
		symbol, err := rand.Int(rand.Reader, big.NewInt(256))
		if err != nil {
			return nil, err
		}

		salt[i] = byte(symbol.Int64())
	}

	return salt, nil
}
