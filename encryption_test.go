package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var original = "I love }apples{"
var samplePasswords = []struct {
	Description string
	Password    string
}{
	{"all_numeric", `01234`},
	{"alpha_lower", `abcdefg`},
	{"alpha_upper", `ABCDEFG`},
	{"empty_string", ``},
	{"single_letter", `a`},
	{"single_number", `5`},
	{"special_chars", `~!@#$%^&*()_+-={}[]\\|;':\",./<>?`},
}

func TestDecryptionReturnsOriginal(t *testing.T) {
	for _, tt := range samplePasswords {
		encrypted := encrypt(original, tt.Password)
		decrypted, _ := decrypt(encrypted, tt.Password)
		msg := fmt.Sprintf("Decryption should return the original (%s)", tt.Description)
		fmt.Printf("Testing password description=%s\n", tt.Description)
		assert.Equal(t, decrypted, original, msg)
	}
	password := "my secret password"
	encrypted := encrypt(original, password)
	decrypted, _ := decrypt(encrypted, password)
	assert.Equal(t, decrypted, original, "Decryption should return the original")

}
