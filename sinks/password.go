package sinks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"go.uber.org/zap"
	"io"
)

type PasswordService interface {
	EncodePassword(plainText string) string
	SetKey(newKey string)
	GetPassword(cipheredText string) string
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

func (ps *passwordService) EncodePassword(plainText string) string {
	cipherText := encrypt([]byte(plainText), ps.key)
	return cipherText
}

func (ps *passwordService) SetKey(newKey string) {
	ps.key = newKey
}

func (ps *passwordService) GetPassword(cipheredText string) string {
	hexedByte, err := hex.DecodeString(cipheredText)
	if err != nil {
		ps.logger.Error("invalid decryption", zap.Error(err))
	}
	plainByte := decrypt(hexedByte, ps.key)

	return string(plainByte)
}

func encrypt(data []byte, passphrase string) string {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext)
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
