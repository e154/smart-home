package encryptor

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

var additionalData = []byte("SMART-HOME")

const (
	//
	// keyLength is the exact key length accepted by Encryptor
	//
	keyLength = 32
	nonceLen  = 24
	tagLen    = 16
)

// Encryptor contains info to encrypt/decrypt sensitive data
type Encryptor struct {
	key []byte
}

// New creates Encryptor using key.
func New(key []byte) (*Encryptor, error) {

	if len(key) != keyLength {
		return nil, errors.New("key must be exactly 32 bytes")
	}

	return &Encryptor{
		key: key,
	}, nil
}

// Encrypt performs XChacha20Poly1305 encryption using saved key.
func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {

	nonce := make([]byte, nonceLen)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	cipher, err := chacha20poly1305.NewX(e.key)
	if err != nil {
		return nil, err
	}
	ct := make([]byte, nonceLen, nonceLen+len(data)+cipher.Overhead())
	copy(ct, nonce)

	return cipher.Seal(ct[:nonceLen], nonce, data, additionalData), nil
}

// Decrypt performs XChacha20Poly1305 decryption using saved key.
func (e *Encryptor) Decrypt(ciphertext []byte) ([]byte, error) {

	if len(ciphertext) < (nonceLen + tagLen) {
		return nil, errors.New("invalid ciphertext length")
	}

	nonce := ciphertext[:nonceLen]

	cipher, err := chacha20poly1305.NewX(e.key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, 0, len(ciphertext)-nonceLen-tagLen)
	return cipher.Open(dst, nonce, ciphertext[nonceLen:], additionalData)

}

var key []byte

func SetKey(value []byte) {
	key = value
}

func GenKey() []byte {
	arr := make([]byte, keyLength)
	rand.Read(arr)
	return arr
}

func Encrypt(value string) (string, error) {
	encryptor, err := New(key)
	if err != nil {
		return "", err
	}
	b, err := encryptor.Encrypt([]byte(value))
	return hex.EncodeToString(b), err
}

func Decrypt(value string) (string, error) {
	encryptor, err := New(key)
	if err != nil {
		return "", err
	}
	b, _ := hex.DecodeString(value)
	decr, err := encryptor.Decrypt(b)
	return string(decr), err
}
