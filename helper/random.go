package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Handle error properly in real applications
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n] // Trim to required length
}

func RandomStringNumber(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Handle error properly in real applications
	}
	var res string
	for i := 0; i < len(b); i++ {
		res += fmt.Sprintf("%d", b[i]%10)
	}
	return res
}
