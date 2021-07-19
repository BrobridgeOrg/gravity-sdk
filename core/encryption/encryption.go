package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Encryption struct {
	Key string
}

func NewEncryption() *Encryption {
	return &Encryption{}
}

func (encryption *Encryption) Encrypt(data []byte) ([]byte, error) {

	block, err := aes.NewCipher([]byte(encryption.Key))
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

	block, err := aes.NewCipher([]byte(encryption.Key))
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
