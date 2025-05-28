package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func ComparePassword(hashedPassword string, plainPassword string) bool {
	return hashedPassword == HashPassword(plainPassword)
}

func HashPassword(plainPassword string) string {
	hash := sha256.New()
	hash.Write([]byte(plainPassword))
	return hex.EncodeToString(hash.Sum(nil))
}
