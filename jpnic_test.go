package jpnic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var caFilePath = "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer"

// HomeNOC
//var certFilePathV4 = "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem"
//var keyFilePathV4 = "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem"
//var certFilePathV6 = "/home/yonedayuto/Documents/HomeNOC/cert/v6-cert.pem"
//var keyFilePathV6 = "/home/yonedayuto/Documents/HomeNOC/cert/v6-prvkey.pem"

// doornoc
var certFilePathV4 = "/home/yonedayuto/Documents/doornoc/cert/v4-cert.pem"
var keyFilePathV4 = "/home/yonedayuto/Documents/doornoc/cert/v4-prvkey.pem"
var certFilePathV6 = "/home/yonedayuto/Documents/doornoc/cert/v6-cert.pem"
var keyFilePathV6 = "/home/yonedayuto/Documents/doornoc/cert/v6-prvkey.pem"

// Search String (HOMENOC/DOORNOC)
var searchStr = "HOMENOC"
var v4UserURL = "/jpnic/entryinfo_v4.do?netwrk_id=2020021426910"
var v6UserURL = "/jpnic/G11320.do?netwrk_id=2020021427992"

var JPNICHandle1 = "YY38053JP"
var JPNICHandle2 = "YY36773JP"

func TestSend(t *testing.T) {
	raw, err := ioutil.ReadFile("./user.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var w WebTransaction

	json.Unmarshal(raw, &w)

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/webtrans/WebRegisterCtl",
		CertFilePath: certFilePathV4,
		KeyFilePath:  keyFilePathV4,
		CAFilePath:   caFilePath,
	}

	result := con.Send(w)
	if result.Err != nil {
		t.Log(result.Err)
		t.Log(result.ResultErr)
		return
	}

	t.Log(result)

	t.Log("受付番号: " + result.RecepNo)
	t.Log("管理者連絡窓口: " + result.AdmJPNICHdl)
	t.Log("技術者連絡窓口1: " + result.Tech1JPNICHdl)
	t.Log("技術者連絡窓口2: " + result.Tech2JPNICHdl)
}

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

	data, err := con.GetAllIPv6(searchStr)
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

func TestGetJPNICHandle1(t *testing.T) {
	con := Config{
		CertFilePath: certFilePathV6,
		KeyFilePath:  keyFilePathV6,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetJPNICHandle(JPNICHandle1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestGetJPNICHandle2(t *testing.T) {
	con := Config{
		CertFilePath: certFilePathV6,
		KeyFilePath:  keyFilePathV6,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetJPNICHandle(JPNICHandle2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestReturnIPv4(t *testing.T) {
	con := Config{
		CertFilePath: certFilePathV4,
		KeyFilePath:  keyFilePathV4,
		CAFilePath:   caFilePath,
	}

	data, err := con.ReturnIPv4("", "Y-NET", "2021/8/31", "noc@doornoc.net")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("受付番号: " + data)
}

func TestReturnIPv6(t *testing.T) {
	con := Config{
		CertFilePath: certFilePathV6,
		KeyFilePath:  keyFilePathV6,
		CAFilePath:   caFilePath,
	}

	// IPv6アドレスの表記はJPNIC側に合わせる必要があります。
	data, err := con.ReturnIPv6([]string{"2407:a2c0:0003::/64"}, "noc@doornoc.net", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("受付番号: " + data)
}
