package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// Reference: https://golang.org/src/crypto/cipher/example_test.go

var (
	// Padding right side padding for CBC to make size uniform for encryption
	Padding = string(rune(0))
	// Bits we're using 16bit IV and 32bit AES
	Bits = struct {
		AES int
		IV  int
	}{32, 16}
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Decrypt a string given the provided password
func Decrypt(content string, password string) (string, error) {
	// Base64 decode the string
	decoded, err := base64.StdEncoding.DecodeString(content)
	handleError(err)

	// Extract the iv and the encrypted cipher
	iv := decoded[:Bits.IV]
	payload := decoded[Bits.IV:]

	// Create the aes thingy
	aesBlock, err := aes.NewCipher(PasswordHash(password))
	handleError(err)

	// Create the decryptor and decrypt payload with it
	decrypter := cipher.NewCBCDecrypter(aesBlock, iv)
	decrypter.CryptBlocks(payload, payload)

	return strings.TrimRight(string(payload), Padding), nil
}

// Encrypt a string given the provided password
func Encrypt(content string, password string) string {
	key := PasswordHash(password)
	block, err := aes.NewCipher(key)
	handleError(err)
	padded := Pad(content)
	ciphertext := make([]byte, aes.BlockSize+len(padded))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(padded))
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Pad with 32bits for AES
func Pad(content string) string {
	return content + strings.Repeat(Padding,
		Bits.AES-len(content)%Bits.AES)
}

// PasswordHash calculates the sha256 hash of a given password
func PasswordHash(password string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}
