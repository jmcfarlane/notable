package app

import (
	"crypto/sha256"
	"encoding/base64"
	"regexp"

	"github.com/gtank/cryptopasta"
)

// Reference: https://golang.org/src/crypto/cipher/example_test.go

var (
	// Decrypted strings generally look like this
	decryptedRE = regexp.MustCompile(`[\x00-\x7F]`)
	// Encrypted string generally look like this
	encryptedRE = regexp.MustCompile(`[^\x00-\x7F]`)
)

// Decrypt a string given the provided password
func Decrypt(content string, password string) (string, error) {
	// Base64 decode the string
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(password))
	cleartext, err := cryptopasta.Decrypt(decoded, &sum)
	if err != nil {
		return "", err
	}
	return string(cleartext), nil
}

// Encrypt a string given the provided password
func Encrypt(content string, password string) (string, error) {

	sum := sha256.Sum256([]byte(password))
	ciphertext, err := cryptopasta.Encrypt([]byte(content), &sum)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// SmellsEncrypted - Try to guess if a string is encrypted or not
func SmellsEncrypted(content string) bool {
	decrypted := len(decryptedRE.FindAllString(content, -1))
	encrypted := len(encryptedRE.FindAllString(content, -1))
	if float64(encrypted)/float64(decrypted) > 0.4 {
		return true
	}
	return false
}
