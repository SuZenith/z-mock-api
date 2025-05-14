package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net/url"
)

func Encryption128Cbc(key, iv, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("create cipher error: %w", err)
	}
	plaintext := pkcs7Padding([]byte(text), block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	b64 := base64.StdEncoding.EncodeToString(ciphertext)

	return url.QueryEscape(b64), nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
