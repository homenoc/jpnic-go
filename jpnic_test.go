package jpnic

import (
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

// HomeNOC
//var certFilePath = "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem"
//var keyFilePath = "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem"
//var caFilePath = "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer"

// doornoc
var certFilePath = "/home/yonedayuto/Documents/doornoc/cert/v4-cert.pem"
var keyFilePath = "/home/yonedayuto/Documents/doornoc/cert/v4-prvkey.pem"
var caFilePath = "/home/yonedayuto/Documents/doornoc/cert/rootcacert_r3.cer"

// Search String
var searchStr = "doornoc"

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
//		CertFilePath: certFilePath,
//		KeyFilePath:  keyFilePath,
//		CAFilePath:   caFilePath,
//	}
//	input := WebTransaction{}
//
//	err := con.Send(input)
//	if err.Err != nil {
//		t.Log(err)
//	}
//}

func Test1GetIPv4(t *testing.T) {
	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetAllIPv4(searchStr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestGetIPv4(t *testing.T) {
	sessionID, err := randomStr()
	if err != nil {
		t.Fatal(err)
	}

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,
		CAFilePath:   caFilePath,
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

	sessionID = resp.Header.Get("Set-Cookie")[11:43]

	cookies = []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	jar.SetCookies(urlObj, cookies)

	resp, err = client.Get("https://iphostmaster.nic.ad.jp/jpnic/portalv4list.do")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	contentType := "application/x-www-form-urlencoded"

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=" + searchStr + "&recepNo=&deliNo="
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

	doc.Find("table").Each(func(_ int, tableHtml *goquery.Selection) {
		tableHtml.Find("tr").Each(func(_ int, rowHtml *goquery.Selection) {
			var info InfoIPv4
			rowHtml.Find("td").Each(func(index int, tableCell *goquery.Selection) {
				ioStr := strings.NewReader(tableCell.Text())
				reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
				strByte, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}
				dataStr := strings.TrimSpace(string(strByte))

				switch index {
				case 0:
					info.IPAddress = dataStr
					info.DetailLink, _ = tableCell.Find("a").Attr("href")
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

func Test1GetIPv6(t *testing.T) {
	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,
		CAFilePath:   caFilePath,
	}

	data, err := con.GetAllIPv6(searchStr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestGetIPv6(t *testing.T) {
	sessionID, err := randomStr()
	if err != nil {
		t.Fatal(err)
	}

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,
		CAFilePath:   caFilePath,
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

	sessionID = resp.Header.Get("Set-Cookie")[11:43]

	cookies = []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	jar.SetCookies(urlObj, cookies)

	contentType := "application/x-www-form-urlencoded"

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=" + searchStr + "&recepNo=&deliNo="
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

	doc.Find("table").Each(func(_ int, tableHtml *goquery.Selection) {
		tableHtml.Find("tr").Each(func(_ int, rowHtml *goquery.Selection) {
			var info InfoIPv6
			rowHtml.Find("td").Each(func(index int, tableCell *goquery.Selection) {
				ioStr := strings.NewReader(tableCell.Text())
				reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
				strByte, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}
				dataStr := strings.TrimSpace(string(strByte))

				switch index {
				case 0:
					info.IPAddress = dataStr
					info.DetailLink, _ = tableCell.Find("a").Attr("href")
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

func TestGetIPv4User(t *testing.T) {
	sessionID, err := randomStr()
	if err != nil {
		t.Fatal(err)
	}

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,
		CAFilePath:   caFilePath,
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

	sessionID = resp.Header.Get("Set-Cookie")[11:43]

	cookies = []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	jar.SetCookies(urlObj, cookies)

	resp, err = client.Get("https://iphostmaster.nic.ad.jp/jpnic/entryinfo_v4.do?netwrk_id=2020021426910")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyStr)))
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	var info InfoDetailIPv4

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						rowHtml2.Find("td").Each(func(_ int, tableCell2 *goquery.Selection) {
							tableCell2.Find("table").Each(func(_ int, tableHtml3 *goquery.Selection) {
								tableHtml3.Find("tr").Each(func(_ int, rowHtml3 *goquery.Selection) {
									rowHtml3.Find("td").Each(func(_ int, tableCell3 *goquery.Selection) {
										tableCell3.Find("table").Each(func(_ int, tableHtml4 *goquery.Selection) {
											tableHtml4.Find("tr").Each(func(_ int, rowHtml4 *goquery.Selection) {
												rowHtml4.Find("td").Each(func(index int, tableCell4 *goquery.Selection) {
													ioStr := strings.NewReader(tableCell4.Text())
													reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
													strByte, err := ioutil.ReadAll(reader)
													if err != nil {
														t.Fatal(err)
													}
													dataStr := strings.TrimSpace(string(strByte))

													if index == 1 {
														switch count {
														case 0:
															info.IPAddress = dataStr
														case 1:
															info.Ryakusho = dataStr
														case 2:
															info.Type = dataStr
														case 3:
															info.InfraUserKind = dataStr
														case 4:
															info.NetworkName = dataStr
														case 5:
															info.Org = dataStr
														case 6:
															info.OrgEn = dataStr
														case 7:
															info.PostCode = dataStr
														case 8:
															info.Address = dataStr
														case 9:
															info.AddressEn = dataStr
														case 10:
															info.AdminJPNICHandle = dataStr
															info.AdminJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case 11:
															info.TechJPNICHandle = dataStr
															info.TechJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case 12:
															info.NameServer = dataStr
														case 13:
															info.DSRecord = dataStr
														case 14:
															info.NotifyAddress = dataStr
														case 15:
															info.DeliNo = dataStr
														case 16:
															info.RecepNo = dataStr
														case 17:
															info.AssignDate = dataStr
														case 18:
															info.ReturnDate = dataStr
														case 19:
															info.UpdateDate = dataStr
														}
														count++
													}
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})

	//infos = infos[1:]

	//for _, tmp := range infos {
	t.Log(info)
	//}
}

func TestGetIPv6User(t *testing.T) {
	sessionID, err := randomStr()
	if err != nil {
		t.Fatal(err)
	}

	con := Config{
		URL:          "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,
		CAFilePath:   caFilePath,
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

	sessionID = resp.Header.Get("Set-Cookie")[11:43]

	cookies = []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	jar.SetCookies(urlObj, cookies)

	resp, err = client.Get("https://iphostmaster.nic.ad.jp/jpnic/G11320.do?netwrk_id=2020021427992")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyStr)))
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	var info InfoDetailIPv6

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						rowHtml2.Find("td").Each(func(_ int, tableCell2 *goquery.Selection) {
							tableCell2.Find("table").Each(func(_ int, tableHtml3 *goquery.Selection) {
								tableHtml3.Find("tr").Each(func(_ int, rowHtml3 *goquery.Selection) {
									rowHtml3.Find("td").Each(func(_ int, tableCell3 *goquery.Selection) {
										tableCell3.Find("table").Each(func(_ int, tableHtml4 *goquery.Selection) {
											tableHtml4.Find("tr").Each(func(_ int, rowHtml4 *goquery.Selection) {
												rowHtml4.Find("td").Each(func(index int, tableCell4 *goquery.Selection) {
													ioStr := strings.NewReader(tableCell4.Text())
													reader := transform.NewReader(ioStr, japanese.ShiftJIS.NewDecoder())
													strByte, err := ioutil.ReadAll(reader)
													if err != nil {
														t.Fatal(err)
													}
													dataStr := strings.TrimSpace(string(strByte))

													if index == 1 {
														switch count {
														case 0:
															info.IPAddress = dataStr
														case 1:
															info.Ryakusho = dataStr
														case 2:
															info.Type = dataStr
														case 3:
															info.InfraUserKind = dataStr
														case 4:
															info.NetworkName = dataStr
														case 5:
															info.Org = dataStr
														case 6:
															info.OrgEn = dataStr
														case 7:
															info.AdminJPNICHandle = dataStr
															info.AdminJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case 8:
															info.TechJPNICHandle = dataStr
															info.TechJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case 9:
															info.AssignDate = dataStr
														case 10:
															info.ReturnDate = dataStr
														case 11:
															info.UpdateDate = dataStr
														}
														count++
													}
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})

	//infos = infos[1:]

	//for _, tmp := range infos {
	t.Log(info)
	//}
}
