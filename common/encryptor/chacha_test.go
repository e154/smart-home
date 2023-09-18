package encryptor

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCreateEncryptorWithTooShortKey try to create encryptor with a too short key.
func TestCreateEncryptorWithTooShortKey(t *testing.T) {

	key := make([]byte, keyLength-1)
	_, err := rand.Read(key)
	require.NoError(t, err)
	_, err = New(key)
	require.Error(t, err)
}

// TestCreateEncryptorWithTooLongKey try to create encryptor with a too long key.
func TestCreateEncryptorWithTooLongKey(t *testing.T) {

	key := make([]byte, keyLength+1)
	_, err := rand.Read(key)
	require.NoError(t, err)
	_, err = New(key)
	require.Error(t, err)
}

// TestEncryptorEncryptDecrypt tests encryption and then correct decryption.
func TestEncryptorEncryptDecrypt(t *testing.T) {

	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	require.NoError(t, err)
	encryptor, err := New(key)
	require.NoError(t, err)

	hexbytes := make([]byte, 256)
	_, err = rand.Read(hexbytes)
	require.NoError(t, err)
	valueStr := hex.EncodeToString(hexbytes)
	valbytes := []byte(valueStr)

	encrypted, err := encryptor.Encrypt(valbytes)
	require.NoError(t, err)
	require.False(t, bytes.Equal(valbytes, encrypted))

	decrypted, err := encryptor.Decrypt(encrypted)
	require.NoError(t, err)
	require.False(t, bytes.Equal(encrypted, decrypted))
	require.True(t, bytes.Equal(valbytes, decrypted))
}

func TestEncryptorEncryptDecryptString(t *testing.T) {
	key = make([]byte, keyLength)
	_, err := rand.Read(key)
	require.NoError(t, err)

	data, err := Encrypt("foobar")
	require.NoError(t, err)

	data, err = Decrypt(data)
	require.NoError(t, err)
	require.Equal(t, data, "foobar")
}

func BenchmarkEncryptor_Encrypt(t *testing.B) {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	require.NoError(t, err)

	data := make([]byte, 16)
	_, err = rand.Read(data)
	require.NoError(t, err)

	encryptor, err := New(key)
	require.NoError(t, err)

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		_, err := encryptor.Encrypt(data)
		require.NoError(t, err)
	}
}

func BenchmarkEncryptor_Decrypt(t *testing.B) {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	require.NoError(t, err)

	data := make([]byte, 16)
	_, err = rand.Read(data)
	require.NoError(t, err)

	encryptor, err := New(key)
	require.NoError(t, err)

	ct, err := encryptor.Encrypt(data)
	require.NoError(t, err)

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {

		pt, err := encryptor.Decrypt(ct)
		require.NoError(t, err)
		require.Equal(t, pt, data)
	}
}
