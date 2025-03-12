package helper

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Handle error properly in real applications
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n] // Trim to required length
}
