package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func decryptFile(encryptedData []byte, key []byte) ([]byte, error) {

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Separate nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	// Extract the nonce from the beginning of the encrypted data
	nonce := encryptedData[:nonceSize]
	// Extract the ciphertext after the nonce
	cipherText := encryptedData[nonceSize:]

	// Decrypt the data
	return gcm.Open(nil, nonce, cipherText, nil)
}
