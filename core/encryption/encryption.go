package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

var tokenPayload = []byte("Brobridge")

type Encryption struct {
	enabled bool
	key     []byte
}

func NewEncryption() *Encryption {
	return &Encryption{}
}

func (encryption *Encryption) GetKey() []byte {
	return encryption.key
}

func (encryption *Encryption) SetKey(key string) {
	if len(key) == 0 {
		encryption.key = []byte("")
		encryption.enabled = false
		return
	}

	hash := sha256.Sum256([]byte(key))
	encryption.key = hash[:]
	encryption.enabled = true
}

func (encryption *Encryption) Encrypt(data []byte) ([]byte, error) {

	if !encryption.enabled {
		return data, nil
	}

	block, err := aes.NewCipher(encryption.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func (encryption *Encryption) Decrypt(data []byte) ([]byte, error) {

	if !encryption.enabled {
		return data, nil
	}

	block, err := aes.NewCipher([]byte(encryption.key))
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (encryption *Encryption) PrepareToken() ([]byte, error) {
	return encryption.Encrypt(tokenPayload)
}

func (encryption *Encryption) ValidateToken(token []byte) bool {
	data, err := encryption.Decrypt(token)
	if err != nil {
		return false
	}

	if bytes.Compare(data, tokenPayload) != 0 {
		return false
	}

	return true
}
