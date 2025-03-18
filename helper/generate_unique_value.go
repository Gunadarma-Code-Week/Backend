package helper

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

func GenerateUniqueFile(originalFileName string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 15)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	// Menambahkan datetime dan huruf acak ke nama file

	millisecond := time.Now().UnixNano() / 1e6

	timeStamp := fmt.Sprintf("%d", millisecond)
	extension := filepath.Ext(originalFileName)
	uniqueFileName := fmt.Sprintf("%s_%s%s", timeStamp, string(b), extension)

	return uniqueFileName
}
