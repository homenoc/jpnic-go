package jpnic

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
)

//func TestResultProcess(t *testing.T) {
//	var result Result
//
//	str := "<html>\n<body>\nRET=00\nRET_CODE=00000000\nRECEP_NO=001\nRECEP_HSNO=0\nADM_JPNIC_HDL=JP00\nADM_GNAME_JP=Y-Net\nTECH1_JPNIC_HDL=JP00\nTECH1_GNAME_JP=Y-Net\nTECH2_JPNIC_HDL=\nTECH2_GNAME_JP=\nCOUNT=0\n</body>\n</html>"
//	reader := strings.NewReader(str)
//	scanner := bufio.NewScanner(reader)
//
//	success := false
//	var retCode []string
//	ret := "00"
//
//	for scanner.Scan() {
//		// RET
//		if strings.Contains(scanner.Text(), "RET=") {
//			t.Logf("Error: %s", scanner.Text()[4:])
//			ret = scanner.Text()[4:]
//			if scanner.Text()[4:] == "00" {
//				success = true
//			}
//			ret = scanner.Text()[4:]
//		}
//
//		// RET_CODE
//		if strings.Contains(scanner.Text(), "RET_CODE=") {
//			t.Logf("RET_CODE Error: %s", scanner.Text()[9:])
//			retCode = append(retCode, scanner.Text()[9:])
//		}
//
//		// RECEP_NO
//		if strings.Contains(scanner.Text(), "RECEP_NO=") {
//			t.Logf("RECEP_NO=%s", scanner.Text()[9:])
//			result.RecepNo = scanner.Text()[9:]
//		}
//
//		// Admin
//		if strings.Contains(scanner.Text(), "ADM_JPNIC_HDL=") {
//			t.Logf("ADM_JPNIC_HDL=%s", scanner.Text()[14:])
//			result.AdmJPNICHdl = scanner.Text()[14:]
//		}
//
//		// Tech1
//		if strings.Contains(scanner.Text(), "TECH1_JPNIC_HDL=") {
//			t.Logf("TECH1_JPNIC_HDL=%s", scanner.Text()[16:])
//			result.Tech1JPNICHdl = scanner.Text()[16:]
//		}
//
//		// Tech2
//		if strings.Contains(scanner.Text(), "TECH2_JPNIC_HDL=") {
//			t.Logf("TECH2_JPNIC_HDL=%s", scanner.Text()[16:])
//			result.Tech2JPNICHdl = scanner.Text()[16:]
//		}
//
//		t.Log(scanner.Text())
//	}
//
//	// RET
//	if ret != "00" {
//		code, _ := strconv.Atoi(ret)
//		ErrorStatusText(code)
//	}
//
//	// RET_CODE
//	for _, code := range retCode {
//		var errStr string
//		t.Logf("%s", code[4:7])
//
//		// interface
//		if code[4:7] == "000" {
//			code, _ := strconv.Atoi(code[4:7])
//			errStr += ErrorStatusText(code)
//
//		}
//
//		// error genre
//		if code[7:] != "0" {
//			code, _ := strconv.Atoi(code[7:])
//			errStr += "_" + ErrorStatusText(code)
//		}
//	}
//	t.Log(success)
//}

//func TestSend(t *testing.T) {
//	con := Config{
//		URL:          "https://iphostmaster.nic.ad.jp/webtrans/WebRegisterCtl",
//		CertFilePath: "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem",
//		KeyFilePath:  "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem",
//		CAFilePath:   "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer",
//	}
//	input := WebTransaction{}
//
//	err := con.Send(input)
//	if err.Err != nil {
//		t.Log(err)
//	}
//}

func TestGetIPv4(t *testing.T) {
	sessionID, err := randomStr()
	if err != nil {
		t.Fatal(err)
	}

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem",
		KeyFilePath:  "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem",
		CAFilePath:   "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer",
	}

	cert, err := tls.LoadX509KeyPair(con.CertFilePath, con.KeyFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Load CA
	caCert, err := ioutil.ReadFile(con.CAFilePath)
	if err != nil {
		t.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	cookies := []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	urlObj, _ := url.Parse("https://iphostmaster.nic.ad.jp/")
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	jar.SetCookies(urlObj, cookies)

	client := &http.Client{Transport: transport, Jar: jar}

	resp, err := client.Get(con.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log(resp.Header)
	t.Log("SESSION_ID:" + resp.Header.Get("Set-Cookie")[11:43])

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		ioStr := strings.NewReader(scanner.Text())
		reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
		strByte, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%s", string(strByte))
	}

	sessionID = resp.Header.Get("Set-Cookie")[11:43]

	cookies = []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	jar.SetCookies(urlObj, cookies)

	//resp, err = client.Get("https://iphostmaster.nic.ad.jp/jpnic/membermenu.do")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer resp.Body.Close()
	//
	//t.Log(resp.Header)
	////t.Log("SESSION_ID:" + resp.Header.Get("Set-Cookie")[11:43])
	//
	//scanner = bufio.NewScanner(resp.Body)
	//for scanner.Scan() {
	//	ioStr := strings.NewReader(scanner.Text())
	//	reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
	//	strByte, err := ioutil.ReadAll(reader)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	t.Logf("%s", string(strByte))
	//}

	resp, err = client.Get("https://iphostmaster.nic.ad.jp/jpnic/portalv4list.do")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log(resp.Header)

	scanner = bufio.NewScanner(resp.Body)

	contentType := "application/x-www-form-urlencoded"

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=HOMENOC&recepNo=&deliNo="
	// utf-8 => shift-jis
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	strByte, err := ioutil.ReadAll(rio)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Post("https://iphostmaster.nic.ad.jp/jpnic/portalv4listmain.do", contentType, bytes.NewBuffer(strByte))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log(resp.Header)

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyStr)))
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	var infos []InfoIPv4

	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			var info InfoIPv4
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				//row = append(row, tablecell.Text())
				ioStr := strings.NewReader(tablecell.Text())
				reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
				strByte, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}
				dataStr := strings.TrimSpace(string(strByte))

				switch indexth {
				case 0:
					info.IPAddress = dataStr
					info.DetailLink, _ = tablecell.Find("a").Attr("href")
				case 1:
					info.Size = dataStr
				case 2:
					info.NetworkName = dataStr
				case 3:
					info.AssignDate = dataStr
				case 4:
					info.ReturnDate = dataStr
				case 5:
					info.OrgName = dataStr
				case 6:
					info.Ryakusho = dataStr
				case 7:
					info.RecepNo = dataStr
				case 8:
					info.DeliNo = dataStr
				case 9:
					info.Type = dataStr
				case 10:
					info.KindID = dataStr
				}
				count++
			})
			if count == 11 {
				infos = append(infos, info)
			}
			count = 0
		})
	})

	infos = infos[1:]

	for _, tmp := range infos {
		t.Log(tmp)
	}
}

func TestGetIPv6(t *testing.T) {
	sessionID, err := randomStr()
	if err != nil {
		t.Fatal(err)
	}

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem",
		KeyFilePath:  "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem",
		CAFilePath:   "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer",
	}

	cert, err := tls.LoadX509KeyPair(con.CertFilePath, con.KeyFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Load CA
	caCert, err := ioutil.ReadFile(con.CAFilePath)
	if err != nil {
		t.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	cookies := []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	urlObj, _ := url.Parse("https://iphostmaster.nic.ad.jp/")
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	jar.SetCookies(urlObj, cookies)

	client := &http.Client{Transport: transport, Jar: jar}

	resp, err := client.Get(con.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log(resp.Header)
	t.Log("SESSION_ID:" + resp.Header.Get("Set-Cookie")[11:43])

	sessionID = resp.Header.Get("Set-Cookie")[11:43]

	cookies = []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	jar.SetCookies(urlObj, cookies)

	contentType := "application/x-www-form-urlencoded"

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=HOMENOC&recepNo=&deliNo="
	// utf-8 => shift-jis
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	strByte, err := ioutil.ReadAll(rio)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Post("https://iphostmaster.nic.ad.jp/jpnic/K11310Action.do", contentType, bytes.NewBuffer(strByte))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log(resp.Header)

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyStr)))
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	var infos []InfoIPv6

	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			var info InfoIPv6
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				//row = append(row, tablecell.Text())
				ioStr := strings.NewReader(tablecell.Text())
				reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
				strByte, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}
				dataStr := strings.TrimSpace(string(strByte))

				switch indexth {
				case 0:
					info.IPAddress = dataStr
					info.DetailLink, _ = tablecell.Find("a").Attr("href")
				case 1:
					info.NetworkName = dataStr
				case 2:
					info.AssignDate = dataStr
				case 3:
					info.ReturnDate = dataStr
				case 4:
					info.OrgName = dataStr
				case 5:
					info.Ryakusho = dataStr
				case 6:
					info.RecepNo = dataStr
				case 7:
					info.DeliNo = dataStr
				case 8:
					info.KindID = dataStr
				}
				count++
			})
			if count == 9 {
				infos = append(infos, info)
			}
			count = 0
		})
	})

	infos = infos[1:]

	for _, tmp := range infos {
		t.Log(tmp)
	}
}
