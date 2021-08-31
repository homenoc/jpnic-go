package jpnic

import (
	"crypto/rand"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"strings"
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

func readShiftJIS(readerInput io.Reader) (string, []byte, error) {
	// shift-jis => utf-8
	reader := transform.NewReader(readerInput, japanese.ShiftJIS.NewDecoder())
	strByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", nil, err

	}
	return string(strByte), strByte, nil
}

func toShiftJIS(str string) (string, []byte, error) {
	// utf-8 => shift-jis
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	strByte, err := ioutil.ReadAll(rio)
	if err != nil {
		return "", nil, err

	}
	return string(strByte), strByte, nil
}
