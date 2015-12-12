package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptionReturnsOriginal(t *testing.T) {
	original := "I love }apples{"
	password := "my secret password"
	encrypted := Encrypt(original, password)
	decrypted, _ := Decrypt(encrypted, password)
	assert.Equal(t, decrypted, original, "Decryption should return the original")

}
