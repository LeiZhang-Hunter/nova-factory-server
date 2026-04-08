// Package aes provides small AES-GCM helpers for encrypting and decrypting string payloads.
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

// DecodeKeyString parses a configured AES key string and normalizes it to a valid AES key length.
//
// It first tries base64 decoding. If decoding fails, it uses the raw string bytes directly.
func DecodeKeyString(key string) ([]byte, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("aes key is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(key)
	if err == nil && len(decoded) > 0 {
		return NormalizeKey(decoded), nil
	}

	return NormalizeKey([]byte(key)), nil
}

// EncryptString encrypts plaintext with AES-GCM and returns base64(nonce+ciphertext).
func EncryptString(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(NormalizeKey(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

// DecryptString decrypts base64(nonce+ciphertext) encrypted by AES-GCM.
func DecryptString(key []byte, encrypted string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(NormalizeKey(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("invalid encrypted payload")
	}

	nonce := raw[:gcm.NonceSize()]
	ciphertext := raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

// NormalizeKey normalizes arbitrary key bytes into a valid AES key length.
//
// Short keys are right-padded with zero bytes, and long keys are truncated.
func NormalizeKey(key []byte) []byte {
	switch {
	case len(key) == 0:
		return make([]byte, 16)
	case len(key) <= 16:
		return normalizeKeyToSize(key, 16)
	case len(key) <= 24:
		return normalizeKeyToSize(key, 24)
	default:
		return normalizeKeyToSize(key, 32)
	}
}

// normalizeKeyToSize pads or truncates key bytes to the target AES key size.
func normalizeKeyToSize(key []byte, size int) []byte {
	out := make([]byte, size)
	copy(out, key)
	return out
}
