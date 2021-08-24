package jpnic

import (
	"log"
	"testing"
)

func SendTest(t *testing.T) {
	con := Config{
		URL:          "",
		CertFilePath: "",
		KeyFilePath:  "",
		CAFilePath:   "",
	}
	input := WebTransaction{}

	err := con.Send(input)
	if err != nil {
		log.Println(err)
	}
}
