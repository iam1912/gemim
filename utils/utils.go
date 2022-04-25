package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func GenToken(id int, nickname string, secret string) string {
	message := fmt.Sprintf("%s%s%d", nickname, secret, id)
	messageMAC := macSha256([]byte(message), []byte(secret))
	return fmt.Sprintf("%s%d", base64.StdEncoding.EncodeToString(messageMAC), id)
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}