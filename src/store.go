package main

import "sync"

type FileData struct {
	EncryptedBytes []byte
	Filename       string
	Key            []byte
}

var(
	FileStore = make(map[string]FileData)
	StoreMutex sync.Mutex
)