package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"mime/multipart"
)

func encryptFile(file []byte, fileHeader *multipart.FileHeader) (FileData, error) {

	// Securely generate a random key 32-byte key for AES-256
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return FileData{}, errors.New("error generating random key")
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return FileData{}, errors.New("error generating cipher block")
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return FileData{}, errors.New("error creating GCM")
	}

	// Create a nonce
	nonce := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)

	// Encrypt the file data
	// {nonce || ciphertext}
	encryptedData := gcm.Seal(nonce, nonce, file, nil)
	return FileData{
		EncryptedBytes: encryptedData,
		Filename:       fileHeader.Filename,
		Key:            key,
	}, nil
}
