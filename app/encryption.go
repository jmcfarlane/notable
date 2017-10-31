package app

import "regexp"

var (
	// Decrypted strings generally look like this
	decryptedRE = regexp.MustCompile(`[\x00-\x7F]`)

	// Encrypted string generally look like this
	encryptedRE = regexp.MustCompile(`[^\x00-\x7F]`)
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
