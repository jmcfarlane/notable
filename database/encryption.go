package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
	iv := decoded[:16]
	payload := decoded[16:]

	// Hash the password
	hasher := sha256.New()
	hasher.Write([]byte(password))

	// Create the aes thingy
	aesBlock, err := aes.NewCipher(hasher.Sum(nil))
	handleError(err)

	// Create the decryptor and decrypt payload with it
	decrypter := cipher.NewCBCDecrypter(aesBlock, iv)
	decrypter.CryptBlocks(payload, payload)

	return string(payload), nil
}
