package hashid

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	// ErrInvalidHash is returned when hash cannot be decoded
	ErrInvalidHash = errors.New("invalid hash")

	// ErrInvalidID is returned when ID is invalid
	ErrInvalidID = errors.New("invalid id")

	// Global encoder instance
	defaultEncoder *Encoder
	once           sync.Once
)

// Encoder handles encoding and decoding of IDs
type Encoder struct {
	block  cipher.Block
	prefix string
}

// NewEncoder creates a new encoder with the given salt
func NewEncoder(salt string, prefix string) (*Encoder, error) {
	// Ensure salt is exactly 16/24/32 bytes for AES
	key := padKey(salt, 16)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return &Encoder{
		block:  block,
		prefix: prefix,
	}, nil
}

// padKey pads or truncates the key to the required length
func padKey(key string, length int) []byte {
	keyBytes := []byte(key)
	if len(keyBytes) >= length {
		return keyBytes[:length]
	}
	// Pad with repeated key if too short
	padded := make([]byte, length)
	for i := 0; i < length; i++ {
		padded[i] = keyBytes[i%len(keyBytes)]
	}
	return padded
}

// Encode converts a uint64 ID to a URL-safe hash string
func (e *Encoder) Encode(id uint64) string {
	// Convert ID to bytes (8 bytes)
	plaintext := make([]byte, 16)
	binary.BigEndian.PutUint64(plaintext[:8], id)
	// Add some randomness-like data based on ID for variety
	binary.BigEndian.PutUint64(plaintext[8:], id^0xDEADBEEF12345678)

	// Encrypt
	ciphertext := make([]byte, 16)
	e.block.Encrypt(ciphertext, plaintext)

	// Encode to URL-safe base64
	encoded := base64.RawURLEncoding.EncodeToString(ciphertext)

	// Add prefix if configured
	if e.prefix != "" {
		return e.prefix + "_" + encoded
	}
	return encoded
}

// Decode converts a hash string back to uint64 ID
func (e *Encoder) Decode(hash string) (uint64, error) {
	// Remove prefix if present
	if e.prefix != "" {
		if strings.HasPrefix(hash, e.prefix+"_") {
			hash = strings.TrimPrefix(hash, e.prefix+"_")
		}
	}

	// Decode from base64
	ciphertext, err := base64.RawURLEncoding.DecodeString(hash)
	if err != nil {
		return 0, ErrInvalidHash
	}

	if len(ciphertext) != 16 {
		return 0, ErrInvalidHash
	}

	// Decrypt
	plaintext := make([]byte, 16)
	e.block.Decrypt(plaintext, ciphertext)

	// Extract ID
	id := binary.BigEndian.Uint64(plaintext[:8])

	// Verify checksum
	checksum := binary.BigEndian.Uint64(plaintext[8:])
	if checksum != id^0xDEADBEEF12345678 {
		return 0, ErrInvalidHash
	}

	return id, nil
}

// MustDecode is like Decode but panics on error
func (e *Encoder) MustDecode(hash string) uint64 {
	id, err := e.Decode(hash)
	if err != nil {
		panic(err)
	}
	return id
}

// Init initializes the global encoder
func Init() error {
	var initErr error
	once.Do(func() {
		salt := os.Getenv("HASHID_SALT")
		if salt == "" {
			salt = "karirnusantara_default_salt_2024" // default, should be changed in production
		}
		prefix := os.Getenv("HASHID_PREFIX")
		if prefix == "" {
			prefix = "kn"
		}

		defaultEncoder, initErr = NewEncoder(salt, prefix)
	})
	return initErr
}

// GetEncoder returns the global encoder instance
func GetEncoder() *Encoder {
	if defaultEncoder == nil {
		if err := Init(); err != nil {
			panic("failed to initialize hashid encoder: " + err.Error())
		}
	}
	return defaultEncoder
}

// Encode is a convenience function that uses the default encoder
func Encode(id uint64) string {
	return GetEncoder().Encode(id)
}

// Decode is a convenience function that uses the default encoder
func Decode(hash string) (uint64, error) {
	return GetEncoder().Decode(hash)
}

// EncodeInt is a convenience function for int
func EncodeInt(id int) string {
	return Encode(uint64(id))
}

// DecodeToInt decodes and returns as int
func DecodeToInt(hash string) (int, error) {
	id, err := Decode(hash)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// EncodeMultiple encodes multiple IDs
func EncodeMultiple(ids []uint64) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = Encode(id)
	}
	return result
}

// DecodeMultiple decodes multiple hashes
func DecodeMultiple(hashes []string) ([]uint64, error) {
	result := make([]uint64, len(hashes))
	for i, hash := range hashes {
		id, err := Decode(hash)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hash at index %d: %w", i, err)
		}
		result[i] = id
	}
	return result, nil
}
