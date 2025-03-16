package helper

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateJoinCode() string {
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	code := hex.EncodeToString(bytes)

	// Format menjadi XXXX-XXXX-XXXX-XXXX-XXXX
	return strings.ToUpper(code[:4] + "-" + code[4:8] + "-" + code[8:12] + "-" + code[12:16] + "-" + code[16:])
}
