package jpnic

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	URL          string
	CertFilePath string
	KeyFilePath  string
	CAFilePath   string
}

func (c *Config) Send(input WebTransaction) Result {
	var result Result

	cert, err := tls.LoadX509KeyPair(c.CertFilePath, c.KeyFilePath)
	if err != nil {
		result.Err = err
		return result
	}

	// Load CA
	caCert, err := ioutil.ReadFile(c.CAFilePath)
	if err != nil {
		result.Err = err
		return result
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	contentType := "text/html"

	str, err := Marshal(input)
	if err != nil {
		result.Err = err
		return result
	}

	// utf-8 => shift-jis
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	strByte, err := ioutil.ReadAll(rio)
	if err != nil {
		result.Err = err
		return result
	}

	resp, err := client.Post(c.URL, contentType, bytes.NewBuffer(strByte))
	if err != nil {
		result.Err = err
		return result
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var retCode []string
	ret := "00"

	for scanner.Scan() {
		// RET
		if strings.Contains(scanner.Text(), "RET=") {
			ret = scanner.Text()[4:]
		}

		// RET_CODE
		if strings.Contains(scanner.Text(), "RET_CODE=") {
			retCode = append(retCode, scanner.Text()[9:])
		}

		// RECEP_NO
		if strings.Contains(scanner.Text(), "RECEP_NO=") {
			result.RecepNo = scanner.Text()[9:]
		}

		// Admin
		if strings.Contains(scanner.Text(), "ADM_JPNIC_HDL=") {
			result.AdmJPNICHdl = scanner.Text()[14:]
		}

		// Tech1
		if strings.Contains(scanner.Text(), "TECH1_JPNIC_HDL=") {
			result.Tech1JPNICHdl = scanner.Text()[16:]
		}

		// Tech2
		if strings.Contains(scanner.Text(), "TECH2_JPNIC_HDL=") {
			result.Tech2JPNICHdl = scanner.Text()[16:]
		}

	}

	// RET
	if ret != "00" {
		code, _ := strconv.Atoi(ret)
		result.Err = fmt.Errorf("%s: %s", ret, ErrorStatusText(code))
	}

	// RET_CODE
	var errStr []error
	for _, codeStr := range retCode {
		var tmpStr string

		// interface
		if codeStr[4:7] != "000" {
			code, _ := strconv.Atoi(codeStr[4:7])
			tmpStr = codeStr[4:7] + ": " + ErrorStatusText(code)

		}

		// error genre
		if codeStr[7:] != "0" {
			code, _ := strconv.Atoi(codeStr[7:])
			tmpStr += "_" + ErrorStatusText(code)
		}

		errStr = append(errStr, fmt.Errorf("%s", tmpStr))
	}

	result.ResultErr = errStr

	return result
}

func (c *Config) GetAllIPv4(searchStr string) ([]InfoIPv4, error) {
	sessionID, err := randomStr()
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(c.CertFilePath, c.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// Load CA
	caCert, err := ioutil.ReadFile(c.CAFilePath)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	jar.SetCookies(urlObj, cookies)

	client := &http.Client{Transport: transport, Jar: jar}

	resp, err := client.Get("https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do")
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	contentType := "application/x-www-form-urlencoded"

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=" + searchStr + "&recepNo=&deliNo="
	// utf-8 => shift-jis
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	strByte, err := ioutil.ReadAll(rio)
	if err != nil {
		return nil, err
	}

	resp, err = client.Post("https://iphostmaster.nic.ad.jp/jpnic/portalv4listmain.do", contentType, bytes.NewBuffer(strByte))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	bodyByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		return nil, err
	}

	count := 0
	var infos []InfoIPv4

	doc.Find("table").Each(func(_ int, tableHtml *goquery.Selection) {
		tableHtml.Find("tr").Each(func(_ int, rowHtml *goquery.Selection) {
			var info InfoIPv4
			rowHtml.Find("td").Each(func(index int, tableCell *goquery.Selection) {
				dataStr := strings.TrimSpace(tableCell.Text())

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

	return infos, nil
}

func (c *Config) GetAllIPv6(searchStr string) ([]InfoIPv6, error) {
	sessionID, err := randomStr()
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(c.CertFilePath, c.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// Load CA
	caCert, err := ioutil.ReadFile(c.CAFilePath)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	jar.SetCookies(urlObj, cookies)

	client := &http.Client{Transport: transport, Jar: jar}

	resp, err := client.Get("https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	resp, err = client.Post("https://iphostmaster.nic.ad.jp/jpnic/K11310Action.do", contentType, bytes.NewBuffer(strByte))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	bodyByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		return nil, err
	}

	count := 0
	var infos []InfoIPv6

	doc.Find("table").Each(func(_ int, tableHtml *goquery.Selection) {
		tableHtml.Find("tr").Each(func(_ int, rowHtml *goquery.Selection) {
			var info InfoIPv6
			rowHtml.Find("td").Each(func(index int, tableCell *goquery.Selection) {
				dataStr := strings.TrimSpace(tableCell.Text())

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

	return infos, nil
}

func (c *Config) GetIPUser(userURL string) (InfoDetailLong, error) {
	var info InfoDetailLong
	var infoShort InfoDetailShort

	sessionID, err := randomStr()
	if err != nil {
		return info, err
	}

	cert, err := tls.LoadX509KeyPair(c.CertFilePath, c.KeyFilePath)
	if err != nil {
		return info, err
	}

	// Load CA
	caCert, err := ioutil.ReadFile(c.CAFilePath)
	if err != nil {
		return info, err
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
		return info, err
	}

	jar.SetCookies(urlObj, cookies)

	client := &http.Client{Transport: transport, Jar: jar}

	resp, err := client.Get("https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do")
	if err != nil {
		return info, err
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

	resp, err = client.Get("https://iphostmaster.nic.ad.jp" + userURL)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	bodyByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return info, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		return info, err
	}

	count := 0

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
													dataStr := strings.TrimSpace(tableCell4.Text())
													if index == 1 {
														switch count {
														case 0:
															info.IPAddress = dataStr
															//short
															infoShort.IPAddress = dataStr
														case 1:
															info.Ryakusho = dataStr
															//short
															infoShort.Ryakusho = dataStr
														case 2:
															info.Type = dataStr
															//short
															infoShort.Type = dataStr
														case 3:
															info.InfraUserKind = dataStr
															//short
															infoShort.InfraUserKind = dataStr
														case 4:
															info.NetworkName = dataStr
															//short
															infoShort.NetworkName = dataStr
														case 5:
															info.Org = dataStr
															//short
															infoShort.Org = dataStr
														case 6:
															info.OrgEn = dataStr
															//short
															infoShort.OrgEn = dataStr
														case 7:
															info.PostCode = dataStr
															//short
															infoShort.AdminJPNICHandle = dataStr
															infoShort.AdminJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case 8:
															info.Address = dataStr
															//short
															infoShort.TechJPNICHandle = dataStr
															infoShort.TechJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case 9:
															info.AddressEn = dataStr
															//short
															infoShort.AssignDate = dataStr
														case 10:
															info.AdminJPNICHandle = dataStr
															info.AdminJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
															//short
															infoShort.ReturnDate = dataStr
														case 11:
															info.TechJPNICHandle = dataStr
															info.TechJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
															//short
															infoShort.UpdateDate = dataStr
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
											// 自所属ではないとき、struct数が13個になる
											// この場合は、並び替えを実施する
											if count == 12 {
												info = InfoDetailLong{}
												info.IPAddress = infoShort.IPAddress
												info.Ryakusho = infoShort.Ryakusho
												info.Type = infoShort.Type
												info.InfraUserKind = infoShort.InfraUserKind
												info.NetworkName = infoShort.NetworkName
												info.Org = infoShort.Org
												info.OrgEn = infoShort.OrgEn
												info.AdminJPNICHandle = infoShort.AdminJPNICHandle
												info.AdminJPNICHandleLink = infoShort.AdminJPNICHandleLink
												info.TechJPNICHandle = infoShort.TechJPNICHandle
												info.TechJPNICHandleLink = infoShort.TechJPNICHandleLink
												info.NameServer = infoShort.NameServer
												info.AssignDate = infoShort.AssignDate
												info.ReturnDate = infoShort.ReturnDate
												info.UpdateDate = infoShort.UpdateDate
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

	return info, err
}
