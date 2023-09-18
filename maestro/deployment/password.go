package deployment

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"go.uber.org/zap"
	"io"
)

type EncryptionService interface {
	EncodePassword(plainText string) (string, error)
	DecodePassword(cipheredText string) (string, error)
}

func NewEncryptionService(logger *zap.Logger, key string) EncryptionService {
	ps := &encryptionService{
		logger: logger,
		key:    key,
	}
	return ps
}

var _ EncryptionService = (*encryptionService)(nil)

type encryptionService struct {
	key    string
	logger *zap.Logger
}

func (ps *encryptionService) EncodePassword(plainText string) (string, error) {
	cipherText, err := encrypt([]byte(plainText), ps.key)
	if err != nil {
		ps.logger.Error("failed to encrypt password", zap.Error(err))
		return "", err
	}
	return cipherText, nil
}

func (ps *encryptionService) DecodePassword(cipheredText string) (string, error) {
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
