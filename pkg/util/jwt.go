package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT generates a JWT token using a Base64-encoded private key
func GenerateJWT(base64PrivateKey string, claims jwt.MapClaims) (string, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(base64PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	// Parse the RSA private key
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("failed to parse private key PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create a new token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token using a Base64-encoded public key
func ValidateJWT(tokenString string, base64PublicKey string) (*jwt.Token, jwt.MapClaims, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	// Parse the RSA public key
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, nil, errors.New("failed to parse public key PEM block")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("public key is not an RSA key")
	}

	// Parse the token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return parsedToken, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}
