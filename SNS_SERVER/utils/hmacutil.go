package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CombineSessionDataAndSignature(sessionData []byte, signature string) []byte {
	separator := []byte("|")
	combined := append(sessionData, separator...)
	combined = append(combined, []byte(signature)...)
	return combined
}

func SeparateSessionDataAndSignature(combined []byte) (string, string, error) {
	separator := []byte("|")
	parts := bytes.SplitN(combined, separator, 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid combined data format")
	}
	data := string(parts[0])
	signature := string(parts[1])
	return data, signature, nil
}

func CreateSignature(value []byte, secretKey string) string {
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write(value)
	return hex.EncodeToString(hmac.Sum(nil))
}
