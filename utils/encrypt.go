package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func Hash1(content []byte) string {
	h := sha1.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(content []byte) string {
	return fmt.Sprintf("%x", md5.Sum(content))
}
