package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
)

// Reference: https://golang.org/src/crypto/cipher/example_test.go

var (
	// Right side padding for CBC to make size uniform for encryption
	padding = string(rune(0))
	// Bits we're using 16bit IV and 32bit AES
	bits = struct {
		AES int
		IV  int
	}{32, 16}
)

// CBCDecrypt a string given the provided password
func CBCDecrypt(content string, password string) (string, error) {
	// Base64 decode the string
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	// Extract the iv and the encrypted cipher
	iv := decoded[:bits.IV]
	payload := decoded[bits.IV:]

	// Create the aes thingy
	aesBlock, err := aes.NewCipher(passwordHash(password))
	if err != nil {
		return "", err
	}

	// Create the decryptor and decrypt payload with it
	decrypter := cipher.NewCBCDecrypter(aesBlock, iv)
	decrypter.CryptBlocks(payload, payload)

	return strings.TrimRight(string(payload), padding), nil
}

// CBCEncrypt a string given the provided password
func CBCEncrypt(content string, password string) (string, error) {
	key := passwordHash(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	padded := pad(content)
	ciphertext := make([]byte, aes.BlockSize+len(padded))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(padded))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Pad with 32bits for AES
func pad(content string) string {
	return content + strings.Repeat(padding,
		bits.AES-len(content)%bits.AES)
}

// Calculate the sha256 hash of a given password
func passwordHash(password string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}
