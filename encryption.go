package main

import (
	"crypto/sha256"
	"errors"
	"regexp"

	"github.com/gtank/cryptopasta"
)

var (
	// Decrypted strings generally look like this
	decryptedRE = regexp.MustCompile(`[\x00-\x7F]`)

	// Encrypted string generally look like this
	encryptedRE = regexp.MustCompile(`[^\x00-\x7F]`)
)

const (
	aesGcm = "AES-GCM"
)

// SmellsEncrypted - Try to guess if a string is encrypted or not
func SmellsEncrypted(content string) bool {
	decrypted := len(decryptedRE.FindAllString(content, -1))
	encrypted := len(encryptedRE.FindAllString(content, -1))
	if float64(encrypted)/float64(decrypted) > 0.4 {
		return true
	}
	return false
}

// Decrypt using the correct cipher type
func Decrypt(note Note, password string) (string, error) {
	if note.CipherType == aesGcm {
		key := sha256.Sum256([]byte(password))
		clearText, err := cryptopasta.Decrypt([]byte(note.Content), &key)
		return string(clearText), err
	}
	return CBCDecrypt(note.Content, password)
}

// Encrypt a note using the currently desired mechanism: AES-GCM
func Encrypt(note Note) (string, string, error) {
	if note.Password == "" {
		return "", "", errors.New("Cannot encrypt with an empty password")
	}
	key := sha256.Sum256([]byte(note.Password))
	cipherText, err := cryptopasta.Encrypt([]byte(note.Content), &key)
	return string(cipherText), aesGcm, err
}
