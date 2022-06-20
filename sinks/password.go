package sinks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func Encrypt(key, data []byte) (cipherText []byte, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}
	cipherText = gcm.Seal(nonce, nonce, data, nil)
	return
}

func Decrypt(key, data []byte) (plaintext []byte, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}
	nonce, cipherText := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err = gcm.Open(nil, nonce, cipherText, nil)
	return
}
