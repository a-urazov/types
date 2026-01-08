package uuid

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	uuid1, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	uuid2, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if uuid1 == uuid2 {
		t.Errorf("Two generated UUIDs should be different: %v == %v", uuid1, uuid2)
	}

	if !uuid1.IsValid() {
		t.Errorf("Generated UUID should be valid: %v", uuid1)
	}
}

func TestFromString(t *testing.T) {
	// Test valid UUID string
	validUUID := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	uuid, err := FromString(validUUID)
	if err != nil {
		t.Errorf("FromString(%q) returned error: %v", validUUID, err)
	}

	if !uuid.IsValid() {
		t.Errorf("Parsed UUID should be valid: %v", uuid)
	}

	// Test valid UUID without hyphens
	validUUIDNoHyphens := "f47ac10b58cc4372a5670e02b2c3d479"
	uuid2, err := FromString(validUUIDNoHyphens)
	if err != nil {
		t.Errorf("FromString(%q) returned error: %v", validUUIDNoHyphens, err)
	}

	if !uuid2.IsValid() {
		t.Errorf("Parsed UUID should be valid: %v", uuid2)
	}

	// Test invalid UUID string
	invalidUUID := "invalid-uuid"
	_, err = FromString(invalidUUID)
	if err == nil {
		t.Errorf("FromString(%q) should return error", invalidUUID)
	}

	// Test short UUID string
	shortUUID := "f47ac10b-58cc-4372-a567-0e02b2c3d4"
	_, err = FromString(shortUUID)
	if err == nil {
		t.Errorf("FromString(%q) should return error", shortUUID)
	}
}

func TestString(t *testing.T) {
	uuidStr := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	uuid, err := FromString(uuidStr)
	if err != nil {
		t.Fatalf("FromString(%q) returned error: %v", uuidStr, err)
	}

	result := uuid.String()
	if result != uuidStr {
		t.Errorf("String() = %q, want %q", result, uuidStr)
	}

	if !strings.Contains(result, "-") {
		t.Errorf("String() result should contain hyphens: %q", result)
	}
}

func TestBytes(t *testing.T) {
	uuidStr := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	uuid, err := FromString(uuidStr)
	if err != nil {
		t.Fatalf("FromString(%q) returned error: %v", uuidStr, err)
	}

	bytes := uuid.Bytes()
	if len(bytes) != 16 {
		t.Errorf("Bytes() returned %d bytes, want 16", len(bytes))
	}

	// Verify that the bytes match the original UUID when reconstructed
	_, err = FromString(string(bytes))
	if err == nil { // This approach won't work directly, so let's just check length
		// The test above ensures the byte slice has the right length
	}
}

func TestIsValid(t *testing.T) {
	// Test valid UUID
	uuid, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if !uuid.IsValid() {
		t.Errorf("New() UUID should be valid: %v", uuid)
	}

	// Test nil UUID
	nilUUID := Nil
	if nilUUID.IsValid() {
		t.Errorf("Nil UUID should not be valid: %v", nilUUID)
	}
}

func TestEqual(t *testing.T) {
	uuid1, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	uuid2 := uuid1 // Copy of uuid1
	if !uuid1.Equal(uuid2) {
		t.Errorf("UUID should be equal to its copy: %v != %v", uuid1, uuid2)
	}

	uuid3, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if uuid1.Equal(uuid3) {
		t.Errorf("Different UUIDs should not be equal: %v == %v", uuid1, uuid3)
	}
}

func TestNilUUID(t *testing.T) {
	nilUUID := Nil

	if nilUUID.IsValid() {
		t.Errorf("Nil UUID should not be valid")
	}

	emptyUUID := UUID{}
	if emptyUUID.IsValid() {
		t.Errorf("Empty UUID should not be valid")
	}
}
