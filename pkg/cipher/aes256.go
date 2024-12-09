package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptAES256 encrypts plaintext using AES-256 GCM.
// The key must be 32 bytes long for AES-256.
func EncryptAES256(plaintext string, keyText string) (string, error) {
	key := []byte(keyText)

	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes long")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	nonce := make([]byte, 12) // GCM standard nonce size is 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine nonce and ciphertext for storage/transmission
	result := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptAES256 decrypts a base64-encoded ciphertext using AES-256 GCM.
// The key must be 32 bytes long.
func DecryptAES256(encryptedText string, keyText string) (string, error) {
	key := []byte(keyText)
	if len(key) != 32 {
		return "", errors.New("decryption key must be 32 bytes long")
	}

	// Decode base64 ciphertext
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	if len(data) < 12 {
		return "", errors.New("invalid ciphertext: too short")
	}

	// Split nonce and ciphertext
	nonce, ciphertext := data[:12], data[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return string(plaintext), nil
}
