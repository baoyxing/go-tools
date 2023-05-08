package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func GetRandAk() string {
	patter := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLOM" +
		"NOPQRSTUVWXYZ123456789"
	ak := ""
	for index := 0; index < 16; index++ {
		n := rand.Intn(61)
		ak += patter[n : n+1]
	}
	return ak
}

func GetSecrect(ak string) string {
	signer := ak + "/" + time.Now().String()
	return fmt.Sprintf("%x", md5.Sum([]byte(signer)))
}
