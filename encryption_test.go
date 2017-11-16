package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEncryption(t *testing.T) {
	password := "my secret password"
	content := "Some beer pls"
	note := Note{Content: content, Password: password}
	cipherText, cipherType, err := Encrypt(note)
	assert.Nil(t, err, "Should not be an error calling encrypt")
	assert.Equal(t, cipherType, "AES-GCM", "Default encryption type is wrong")
	note.CipherType = cipherType
	note.Content = cipherText
	decrypted, err := Decrypt(note, password)
	assert.Nil(t, err, "Should not be an error calling decrypt")
	assert.Equal(t, decrypted, content, "Decryption should return the original")
}

func TestDecryptionOfDeprecatedCBC(t *testing.T) {
	password := "my secret password"
	content := "Some beer pls"
	encrypted, err := CBCEncrypt(content, password)
	assert.Nil(t, err, "Should not be an error calling encrypt")
	note := Note{Content: encrypted, Encrypted: true}
	decrypted, err := Decrypt(note, password)
	assert.Nil(t, err, "Should not be an error calling decrypt")
	assert.Equal(t, decrypted, content, "Decryption should return the original")
}

func TestEncryptNoteWithEmptyPassword(t *testing.T) {
	note := Note{Content: "", Password: ""}
	cipherText, cipherType, err := Encrypt(note)
	assert.NotNil(t, err)
	assert.Empty(t, cipherText)
	assert.Empty(t, cipherType)
}

func TestSmellsEncrypted(t *testing.T) {
	assert.True(t, SmellsEncrypted(`Y���JQ3�����/��z#�4+4���X��'���N8u`))
	assert.False(t, SmellsEncrypted(`Hello world`))
}
