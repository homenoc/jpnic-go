package jpnic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var caFilePath = "/Users/y-yoneda/github/homenoc/jpnic-go/cert/rootcacert_r3.cer"

//HomeNOC
var pfxPass = "homenoc"
var pfxFilePathV4 = "/Users/y-yoneda/github/homenoc/jpnic-go/cert/v4-openssl.p12"
var pfxFilePathV6 = "/Users/y-yoneda/github/homenoc/jpnic-go/cert/v6-openssl.p12"

// HomeNOC
//var pfxFilePathV4 = "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem"
//var pfxFilePathV6 = "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem"

// doornoc
//var pfxFilePathV4 = "/home/yonedayuto/Documents/doornoc/cert/v4-cert.pem"
//var pfxFilePathV6 = "/home/yonedayuto/Documents/doornoc/cert/v4-prvkey.pem"

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
		URL:         "https://iphostmaster.nic.ad.jp/webtrans/WebRegisterCtl",
		PfxFilePath: pfxFilePathV4,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
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

func TestSearchIPv4(t *testing.T) {
	con := Config{
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		PfxFilePath: pfxFilePathV4,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	search := SearchIPv4{
		IPAddress:      "",
		SizeStart:      "",
		SizeEnd:        "",
		NetworkName:    "",
		RegStart:       "",
		RegEnd:         "",
		ReturnStart:    "",
		ReturnEnd:      "",
		Org:            "",
		Ryakusho:       "doornoc",
		RecepNo:        "",
		DeliNo:         "",
		IsAllocate:     false,
		IsAssignInfra:  false,
		IsAssignUser:   false,
		IsSubAllocate:  false,
		IsHistoricalPI: false,
		IsSpecialPI:    false,
	}

	data, jpnicHandles, err := con.SearchIPv4(false, true, search)
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data {
		t.Log(tmp)
	}

	for _, tmp := range jpnicHandles {
		t.Log(tmp)
	}
}

func TestSearchOurIPv4(t *testing.T) {
	con := Config{
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		PfxFilePath: pfxFilePathV4,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, jpnicHandles, err := con.SearchIPv4(true, true, SearchIPv4{})
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data {
		t.Log(tmp)
	}

	for _, tmp := range jpnicHandles {
		t.Log(tmp)
	}
}

func TestSearchIPv6(t *testing.T) {
	con := Config{
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		PfxFilePath: pfxFilePathV6,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	search := SearchIPv6{
		IPAddress:     "",
		SizeStart:     "",
		SizeEnd:       "",
		NetworkName:   "",
		RegStart:      "",
		RegEnd:        "",
		ReturnStart:   "",
		ReturnEnd:     "",
		Org:           "",
		Ryakusho:      "doornoc",
		RecepNo:       "",
		DeliNo:        "",
		IsAllocate:    false,
		IsAssignInfra: false,
		IsAssignUser:  false,
		IsSubAllocate: false,
	}

	data, jpnicHandles, err := con.SearchIPv6(false, true, search)
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data {
		t.Log(tmp)
	}

	for _, tmp := range jpnicHandles {
		t.Log(tmp)
	}
}

func TestSearchOurIPv6(t *testing.T) {
	con := Config{
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		PfxFilePath: pfxFilePathV6,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, jpnicHandles, err := con.SearchIPv6(true, true, SearchIPv6{})
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data {
		t.Log(tmp)
	}

	for _, tmp := range jpnicHandles {
		t.Log(tmp)
	}
}

func TestGetIPv4User(t *testing.T) {
	con := Config{
		PfxFilePath: pfxFilePathV4,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, err := con.GetIPUser(v4UserURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)
}

func TestGetIPv6User(t *testing.T) {
	con := Config{
		PfxFilePath: pfxFilePathV6,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, err := con.GetIPUser(v6UserURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestGetJPNICHandle1(t *testing.T) {
	con := Config{
		PfxFilePath: pfxFilePathV6,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, err := con.GetJPNICHandle(JPNICHandle1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestGetJPNICHandle2(t *testing.T) {
	con := Config{
		PfxFilePath: pfxFilePathV6,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, err := con.GetJPNICHandle(JPNICHandle2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

//func TestReturnIPv4(t *testing.T) {
//	con := Config{
//		PfxFilePath: pfxFilePathV4,
//		PfxPass:     pfxPass,
//		CAFilePath:  caFilePath,
//	}
//
//	data, err := con.ReturnIPv4("", "Y-NET", "2021/8/31", "noc@doornoc.net")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("受付番号: " + data)
//}
//
//func TestReturnIPv6(t *testing.T) {
//	con := Config{
//		PfxFilePath: pfxFilePathV6,
//		PfxPass:     pfxPass,
//		CAFilePath:  caFilePath,
//	}
//
//	// IPv6アドレスの表記はJPNIC側に合わせる必要があります。
//	data, err := con.ReturnIPv6([]string{"2407:a2c0:0003::/64"}, "noc@doornoc.net", "")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("受付番号: " + data)
//}
//
func TestChangeUserInfo(t *testing.T) {
	con := Config{
		PfxFilePath: pfxFilePathV4,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	raw, err := ioutil.ReadFile("./user_detail.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var input JPNICHandleInput

	json.Unmarshal(raw, &input)

	data, err := con.ChangeUserInfo(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("受付番号: " + data)
}

//func TestRequestInfo(t *testing.T) {
//	con := Config{
//		PfxFilePath: pfxFilePathV4,
//		PfxPass:     pfxPass,
//		CAFilePath:  caFilePath,
//	}
//
//	data, err := con.GetRequestList("")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for _, tmp := range data {
//		t.Log(tmp)
//	}
//}

//func TestRecepInfo(t *testing.T) {
//	con := Config{
//		PfxFilePath: pfxFilePathV4,
//		PfxPass:     pfxPass,
//		CAFilePath:  caFilePath,
//	}
//
//	data, err := con.GetDetailRequest("020210816000002")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	t.Log(data)
//}

func TestGetResourceManagement(t *testing.T) {
	con := Config{
		PfxFilePath: pfxFilePathV4,
		PfxPass:     pfxPass,
		CAFilePath:  caFilePath,
	}

	data, html, err := con.GetResourceManagement()
	if err != nil {
		t.Fatal(err)
	}

	for _, tmp := range data.ResourceCIDRBlock {
		t.Log(tmp)
	}

	t.Log("--------------HTML--------------")

	t.Log(html)
}
