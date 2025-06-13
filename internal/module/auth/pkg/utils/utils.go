package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
)

func GenerateCode(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)

	for i := range code {
		code[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(code)
}

func GeneratePassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashString(s string) (string, error) {
	h, err := blake2b.New256(nil)
	if err != nil {
		return "", err
	}
	h.Write([]byte(s))
	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

func GenerateSecretToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := crand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil

}
