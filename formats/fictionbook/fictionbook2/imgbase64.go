package fictionbook2

import (
	"encoding/base64"
)

// DecodeBase64 returns the bytes represented by the base64 string s
func DecodeBase64(source string) ([]byte, error) {
	out, err := base64.StdEncoding.DecodeString(source)
	return out, err
}
