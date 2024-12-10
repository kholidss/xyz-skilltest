package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToSha256(input string) string {
	hasher := sha256.New()

	hasher.Write([]byte(input))

	hashBytes := hasher.Sum(nil)

	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
