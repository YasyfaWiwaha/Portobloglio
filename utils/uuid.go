package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
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
