package utils

import "strings"

type Token struct{}

// ExtractBearerToken Remove "Bearer " from "Authorization" token string
func (tk Token) ExtractBearerToken(t string) (string, error) {
	f := strings.Split(t, " ")
	if len(f) != 2 || f[0] != "Bearer" {
		return "", nil
	}
	return f[1], nil
}
