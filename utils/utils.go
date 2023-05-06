package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Token struct{}

// ExtractBearerToken Remove "Bearer " from "Authorization" token string
func (tk Token) ExtractBearerToken(t string) (string, error) {
	f := strings.Split(t, " ")
	if len(f) != 2 || f[0] != "Bearer" {
		return "", nil
	}
	return f[1], nil
}

func GenerateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	num := ""
	for i := 0; i < 10; i++ {
		num += strconv.Itoa(rand.Intn(10))
	}
	return num
}
