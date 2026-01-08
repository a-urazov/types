package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// UUID represents a universally unique identifier
type UUID [16]byte

// Nil is the zero UUID (all zeros)
var Nil = UUID{}

// New generates a new random UUID
func New() (UUID, error) {
	var uuid UUID
	if _, err := rand.Read(uuid[:]); err != nil {
		return Nil, err
	}

	// Set version (4) and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC 4122

	return uuid, nil
}

// FromString creates a UUID from a string representation
func FromString(s string) (UUID, error) {
	var uuid UUID

	// Remove hyphens if present
	s = strings.ReplaceAll(s, "-", "")

	if len(s) != 32 {
		return Nil, errors.New("invalid UUID string length")
	}

	bytes, err := hex.DecodeString(s)
	if err != nil {
		return Nil, fmt.Errorf("invalid UUID string: %w", err)
	}

	if len(bytes) != 16 {
		return Nil, errors.New("invalid UUID byte length")
	}

	copy(uuid[:], bytes)
	return uuid, nil
}

// String returns the string representation of the UUID
func (uuid UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16])
}

// Bytes returns the underlying byte array
func (uuid UUID) Bytes() []byte {
	return uuid[:]
}

// IsValid checks if the UUID is valid (not nil)
func (uuid UUID) IsValid() bool {
	return uuid != Nil
}

// Equal checks if two UUIDs are equal
func (uuid UUID) Equal(other UUID) bool {
	return uuid == other
}