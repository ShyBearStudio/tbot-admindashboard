package dbutils

import (
	"crypto/sha1"
	"fmt"
)

// Hashes plaintext with SHA-1
func Encrypt(textToEncrypt string) (encryptedText string) {
	encryptedText = fmt.Sprintf("%x", sha1.Sum([]byte(textToEncrypt)))
	return
}
