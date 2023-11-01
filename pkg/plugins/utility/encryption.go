package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// NewRandomSymmetricKey creates a new random symmetric key
func NewRandomSymmetricKey() []byte {
	key := make([]byte, 32) // AES-256 bit key
	_, err := rand.Read(key)
	if err != nil {
		panic(err) // handle error appropriately
	}
	return key
}

// DecodeKeyString decodes a base64 encoded string into a byte array
func DecodeKeyString(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(key)
}

// NewRpcEncryptor creates a new RpcEncryptor with the given symmetric key
func NewRpcEncryptor(symmetricKey []byte) *RpcEncryptor {
	return &RpcEncryptor{
		SymmetricKey: symmetricKey,
	}
}

// RpcEncryptor is used to encrypt and decrypt data using a symmetric key
type RpcEncryptor struct {
	SymmetricKey []byte
}

// KeyString returns the base64 encoded symmetric key
func (r *RpcEncryptor) KeyString() string {
	return base64.StdEncoding.EncodeToString(r.SymmetricKey)
}

// Encrypt encrypts the given data using the symmetric key
func (r *RpcEncryptor) Encrypt(data string) (string, error) {
	block, err := aes.NewCipher(r.SymmetricKey)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

// Decrypt decrypts the given data using the symmetric key
func (r *RpcEncryptor) Decrypt(encodedData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(r.SymmetricKey)
	if err != nil {
		return "", err
	}

	if len(data) < 12 {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:12], data[12:]
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
