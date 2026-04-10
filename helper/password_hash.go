package helper

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	h := sha256.New()
	h.Write([]byte(password))
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	bytes, err := bcrypt.GenerateFromPassword([]byte(sha256Hash), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	h := sha256.New()
	h.Write([]byte(password))
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(sha256Hash))
	return err == nil
}
