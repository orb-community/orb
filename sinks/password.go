package sinks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"go.uber.org/zap"
	"io"
)

type PasswordService interface {
	EncodePassword(plainText string) (string, error)
	SetKey(newKey string)
	DecodePassword(cipheredText string) (string, error)
}

func NewPasswordService(logger *zap.Logger, key string) *passwordService {
	ps := &passwordService{
		logger: logger,
	}
	ps.SetKey(key)
	return ps
}

type passwordService struct {
	key    string
	logger *zap.Logger
}

func (ps *passwordService) EncodePassword(plainText string) (string, error) {
	cipherText, err := encrypt([]byte(plainText), ps.key)
	if err != nil {
		ps.logger.Error("failed to encrypt password", zap.Error(err))
		return "", err
	}
	return cipherText, nil
}

func (ps *passwordService) SetKey(newKey string) {
	ps.key = newKey
}

func (ps *passwordService) DecodePassword(cipheredText string) (string, error) {
	hexedByte, err := hex.DecodeString(cipheredText)
	if err != nil {
		ps.logger.Error("failed to decode password", zap.Error(err))
		return "", err
	}
	plainByte, err := decrypt(hexedByte, ps.key)
	if err != nil {
		ps.logger.Error("failed to decrypt password", zap.Error(err))
		return "", err
	}

	return string(plainByte), nil
}

func encrypt(data []byte, passphrase string) (string, error) {
	block, _ := aes.NewCipher(createHash(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := createHash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func createHash(key string) []byte {
	hasher := sha256.Sum256([]byte(key))
	return hasher[:]
}
