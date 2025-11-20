package utility

import (
	sharedutil "github.com/bgrewell/dtac-agent/pkg/shared/utility"
)

// NewRandomSymmetricKey creates a new random symmetric key
func NewRandomSymmetricKey() []byte {
	return sharedutil.NewRandomSymmetricKey()
}

// DecodeKeyString decodes a base64 encoded string into a byte array
func DecodeKeyString(key string) ([]byte, error) {
	return sharedutil.DecodeKeyString(key)
}

// NewRPCEncryptor creates a new RPCEncryptor with the given symmetric key
func NewRPCEncryptor(symmetricKey []byte) *RPCEncryptor {
	return &RPCEncryptor{
		inner: sharedutil.NewRPCEncryptor(symmetricKey),
	}
}

// RPCEncryptor is used to encrypt and decrypt data using a symmetric key
type RPCEncryptor struct {
	inner *sharedutil.RPCEncryptor
}

// KeyString returns the base64 encoded symmetric key
func (r *RPCEncryptor) KeyString() string {
	return r.inner.KeyString()
}

// Encrypt encrypts the given data using the symmetric key
func (r *RPCEncryptor) Encrypt(data string) (string, error) {
	return r.inner.Encrypt(data)
}

// Decrypt decrypts the given data using the symmetric key
func (r *RPCEncryptor) Decrypt(encodedData string) (string, error) {
	return r.inner.Decrypt(encodedData)
}
