package jpnic

import (
	"crypto/rand"
)

func randomStr() (string, error) {
	const str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(str[int(v)%len(str)])
	}
	return result, nil
}
