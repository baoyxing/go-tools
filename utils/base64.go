package utils

import "encoding/base64"

func Base64EncodeToString(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
