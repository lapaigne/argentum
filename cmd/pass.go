package main

import (
	"crypto/rand"
	"math/big"
)

func genPass(length int) (string, error) {
	pass := make([]byte, length)
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"

	for i := range length {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		pass[i] = chars[r.Int64()]
	}

	return string(pass), nil
}
