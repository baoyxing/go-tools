package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash1(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
