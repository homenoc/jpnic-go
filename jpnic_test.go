package jpnic

import (
	"testing"
)

var caFilePath = "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer"

// HomeNOC
var certFilePathV4 = "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem"
var keyFilePathV4 = "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem"
var certFilePathV6 = "/home/yonedayuto/Documents/HomeNOC/cert/v6-cert.pem"
var keyFilePathV6 = "/home/yonedayuto/Documents/HomeNOC/cert/v6-prvkey.pem"

// doornoc
//var certFilePathV4 = "/home/yonedayuto/Documents/doornoc/cert/v4-cert.pem"
//var keyFilePathV4 = "/home/yonedayuto/Documents/doornoc/cert/v4-prvkey.pem"
//var certFilePathV6 = "/home/yonedayuto/Documents/doornoc/cert/v6-cert.pem"
//var keyFilePathV6 = "/home/yonedayuto/Documents/doornoc/cert/v6-prvkey.pem"

// Search String (HOMENOC/DOORNOC)
var searchStr = "HOMENOC"
var v4UserURL = "/jpnic/entryinfo_v4.do?netwrk_id=2020021426910"
var v6UserURL = "/jpnic/G11320.do?netwrk_id=2020021427992"

func TestGetIPv4(t *testing.T) {
	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePathV4,
		KeyFilePath:  keyFilePathV4,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetAllIPv4(searchStr)
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data {
		t.Log(tmp)

	}
}

func TestGetIPv6(t *testing.T) {
	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePathV6,
		KeyFilePath:  keyFilePathV6,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetAllIPv4(searchStr)
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data {
		t.Log(tmp)
	}
}

func TestGetIPv4User(t *testing.T) {
	con := Config{
		CertFilePath: certFilePathV4,
		KeyFilePath:  keyFilePathV4,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetIPUser(v4UserURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)
}

func TestGetIPv6User(t *testing.T) {
	con := Config{
		CertFilePath: certFilePathV6,
		KeyFilePath:  keyFilePathV6,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetIPUser(v6UserURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
