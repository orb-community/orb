package sinks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"go.uber.org/zap"
	"os"
)

type PasswordService interface {
	EncodePassword(plainText string) (string, error)
	SetKey(newKey []byte)
	GetPassword(cipheredText string) string
}

func NewInstance(logger zap.Logger) *passwordServices {
	keyString := os.Getenv("ORB_SINK_SECRET_KEY")
	if keyString != "" {
		logger.Error("not found the ORB SINK SECRET")
		return nil
	}
	ps := &passwordServices{
		logger: logger,
	}
	ps.SetKey([]byte(keyString))
	return ps
}

type passwordServices struct {
	key    []byte
	logger zap.Logger
}

func (ps *passwordServices) EncodePassword(plainText string) (string, error) {
	cipherText, err := encrypt(ps.key, []byte(plainText))
	if err != nil {
		ps.logger.Error("invalid encryption", zap.Error(err))
		return "", err
	}
	return string(cipherText), nil
}

func (ps *passwordServices) SetKey(newKey []byte) {
	blockCipher, err := aes.NewCipher(newKey)
	if err != nil {
		ps.logger.Error("invalid key", zap.Error(err))
		return
	}
	_, err = cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}
	ps.key = newKey
}

func (ps *passwordServices) GetPassword(cipheredText string) string {
	plainByte, err := decrypt(ps.key, []byte(cipheredText))
	if err != nil {
		ps.logger.Error("invalid decryption", zap.Error(err))
	}
	return string(plainByte)
}

func encrypt(key, data []byte) (cipherText []byte, err error) {
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

func decrypt(key, data []byte) (plaintext []byte, err error) {
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
