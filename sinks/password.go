package sinks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"go.uber.org/zap"
)

type PasswordService interface {
	EncodePassword(plainText string) (string, error)
	SetKey(newKey []byte)
	GetPassword(cipheredText string) string
}

func NewPasswordService(logger *zap.Logger, key string) *passwordService {
	ps := &passwordService{
		logger: logger,
	}
	ps.SetKey([]byte(key))
	return ps
}

type passwordService struct {
	key    []byte
	logger *zap.Logger
}

// Gets the Password encrypted and in Base64 string for storing
func (ps *passwordService) EncodePassword(plainText string) (string, error) {
	cipherText, err := encrypt(ps.key, []byte(plainText))
	if err != nil {
		ps.logger.Error("invalid encryption", zap.Error(err))
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (ps *passwordService) SetKey(newKey []byte) {
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

// Gets the Password from the Base64 string we store
func (ps *passwordService) GetPassword(cipheredText string) string {
	var cipheredByte []byte
	_, err := base64.StdEncoding.Decode(cipheredByte, []byte(cipheredText))
	if err != nil {
		ps.logger.Error("invalid decoding", zap.Error(err))
	}
	plainByte, err := decrypt(ps.key, cipheredByte)
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
