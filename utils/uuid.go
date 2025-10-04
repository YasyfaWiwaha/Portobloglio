package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func FormatUUID(b []byte) (string, error) {
	if len(b) != 16 {
		return "", fmt.Errorf("invalid UUID length: %d", len(b))
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]),
		uint16(b[4])<<8|uint16(b[5]),
		uint16(b[6])<<8|uint16(b[7]),
		uint16(b[8])<<8|uint16(b[9]),
		b[10:]), nil
}

func ParseUUID(uuidStr string) ([]byte, error) {
	clean := strings.ReplaceAll(uuidStr, "-", "")
	if len(clean) != 32 {
		return nil, fmt.Errorf("invalid UUID length: %d", len(clean))
	}

	bytes, err := hex.DecodeString(clean)
	if err != nil {
		return nil, fmt.Errorf("failed to decode UUID hex: %w", err)
	}

	return bytes, nil
}

func GenerateUUIDv7() (dbBlob []byte, structID string, err error) {
	dbBlob = make([]byte, 16)

	timestamp := time.Now().UnixMilli()
	binary.BigEndian.PutUint64(dbBlob[0:8], uint64(timestamp))
	copy(dbBlob[0:6], dbBlob[2:8])
	_, err = rand.Read(dbBlob[6:])
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Set version 7 in byte 6 (bits 4-7)
	dbBlob[6] = (dbBlob[6] & 0x0f) | 0x70

	// Set variant bits in byte 8 (bits 6-7 = 10)
	dbBlob[8] = (dbBlob[8] & 0x3f) | 0x80

	// Format as string for struct
	structID, err = FormatUUID(dbBlob)
	if err != nil {
		return nil, "", fmt.Errorf("failed to format UUID: %w", err)
	}

	return dbBlob, structID, nil
}
