package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func generateToken()string{

	//Generate 16 random bytes
	tokenBytes :=make([]byte,16)
	io.ReadFull(rand.Reader,tokenBytes)

	//Convert token bytes to hex string 
	token := hex.EncodeToString(tokenBytes)
	return token
}