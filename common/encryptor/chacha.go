// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package encryptor

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/e154/smart-home/common/logger"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	log = logger.MustGetLogger("encryptor")
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

func EncryptBind(data string) string {
	value, err := Encrypt(data)
	if err != nil {

		log.Error(err.Error())
		return ""
	}
	return value
}

func DecryptBind(data string) string {
	value, err := Decrypt(data)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return value
}
