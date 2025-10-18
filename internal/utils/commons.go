package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/skip2/go-qrcode"
)

func GetAddressQRS(addresses []string) (string, error) {
	// Join addresses into a single string
	data := strings.Join(addresses, ",")
	// Path to temp folder
	tempDir := filepath.Join(constants.ApplicationDataDir, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}
	qrPath := filepath.Join(tempDir, "addresses_qr.png")
	// Generate QR code and write to file
	err := qrcode.WriteFile(data, qrcode.Medium, 256, qrPath)
	if err != nil {
		return "", err
	}
	return qrPath, nil
}

func MarshalPayload(payload interface{}) json.RawMessage {
	b, _ := json.Marshal(payload)
	return b
}
