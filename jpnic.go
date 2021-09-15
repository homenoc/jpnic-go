package jpnic

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:91.0) Gecko/20100101 Firefox/91.0"
var contentType = "application/x-www-form-urlencoded"

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

	//req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

	contentType = "text/html"

	str, err := Marshal(input)
	if err != nil {
		result.Err = err
		return result
	}

	// utf-8 => shift-jis
	_, strByte, err := toShiftJIS(str)
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
	client, err := c.initAccess()
	if err != nil {
		return nil, err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=" +
		searchStr + "&recepNo=&deliNo=" + "&action=%81%40%8C%9F%8D%F5%81%40"
	// utf-8 => shift-jis
	reqBody, _, err := toShiftJIS(str)
	if err != nil {
		return nil, err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/portalv4listmain.do",
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return nil, err
	}

	_, bodyByte, err := readShiftJIS(resp.Body)
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
	client, err := c.initAccess()
	if err != nil {
		return nil, err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	str := "destdisp=D12204&ipaddr=&sizeS=&sizeE=&netwrkName=&regDateS=&regDateE=&rtnDateS=&rtnDateE=&organizationName=&resceAdmSnm=" + searchStr + "&recepNo=&deliNo="
	// utf-8 => shift-jis
	reqBody, _, err := toShiftJIS(str)
	if err != nil {
		return nil, err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/K11310Action.do",
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
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

func (c *Config) GetIPUser(userURL string) (InfoDetail, error) {
	var info InfoDetail

	client, err := c.initAccess()
	if err != nil {
		return info, err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + userURL,
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return info, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return info, err
	}

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
												var title string
												rowHtml4.Find("td").Each(func(index int, tableCell4 *goquery.Selection) {
													dataStr := strings.TrimSpace(tableCell4.Text())
													if index == 1 {
														switch title {
														case "IPネットワークアドレス":
															info.IPAddress = dataStr
														case "資源管理者略称":
															info.Ryakusho = dataStr
														case "アドレス種別":
															info.Type = dataStr
														case "インフラ・ユーザ区分":
															info.InfraUserKind = dataStr
														case "ネットワーク名":
															info.NetworkName = dataStr
														case "組織名":
															info.Org = dataStr
														case "Organization":
															info.OrgEn = dataStr
														case "郵便番号":
															info.PostCode = dataStr
														case "住所":
															info.Address = dataStr
														case "Address":
															info.AddressEn = dataStr
														case "管理者連絡窓口":
															info.AdminJPNICHandle = dataStr
															info.AdminJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case "技術連絡担当者":
															info.TechJPNICHandle = dataStr
															info.TechJPNICHandleLink, _ = tableCell4.Find("a").Attr("href")
														case "ネームサーバ":
															info.NameServer = dataStr
														case "DSレコード":
															info.DSRecord = dataStr
														case "通知アドレス":
															info.NotifyAddress = dataStr
														case "審議番号":
															info.DeliNo = dataStr
														case "受付番号":
															info.RecepNo = dataStr
														case "割当年月日":
															info.AssignDate = dataStr
														case "返却年月日":
															info.ReturnDate = dataStr
														case "最終更新":
															info.UpdateDate = dataStr
														}
													} else {
														title = dataStr
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

	return info, err
}

func (c *Config) GetJPNICHandle(handle string) (JPNICHandleDetail, error) {
	var info JPNICHandleDetail

	client, err := c.initAccess()
	if err != nil {
		return info, err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/entryinfo_handle.do?jpnic_hdl=" + handle,
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return info, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return info, err
	}

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						rowHtml2.Find("td").Each(func(_ int, tableCell2 *goquery.Selection) {
							tableCell2.Find("table").Each(func(_ int, tableHtml3 *goquery.Selection) {
								tableHtml3.Find("tr").Each(func(_ int, rowHtml3 *goquery.Selection) {
									var title string
									rowHtml3.Find("td").Each(func(index int, tableCell3 *goquery.Selection) {
										dataStr := strings.TrimSpace(tableCell3.Text())
										if index == 1 {
											switch title {
											case "グループハンドル":
												info.IsJPNICHandle = false
												info.JPNICHandle = dataStr
											case "グループ名":
												info.Org = dataStr
											case "Group Name":
												info.OrgEn = dataStr
											case "JPNICハンドル":
												info.IsJPNICHandle = true
												info.JPNICHandle = dataStr
											case "氏名":
												info.Org = dataStr
											case "Last, First":
												info.OrgEn = dataStr
											case "電子メール":
												info.Email = dataStr
											case "電子メイル": // JPNIC側の表記ゆれのため
												info.Email = dataStr
											case "組織名":
												info.Org = dataStr
											case "Organization":
												info.OrgEn = dataStr
											case "部署":
												info.Division = dataStr
											case "Division":
												info.DivisionEn = dataStr
											case "肩書":
												info.Title = dataStr
											case "Title":
												info.TitleEn = dataStr
											case "電話番号":
												info.Tel = dataStr
											case "Fax番号":
												info.Fax = dataStr
											case "FAX番号": // JPNIC側の表記ゆれのため
												info.Fax = dataStr
											case "通知アドレス":
												info.NotifyAddress = dataStr
											case "最終更新":
												info.UpdateDate = dataStr
											}
										} else {
											title = dataStr
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

	return info, err
}

func (c *Config) ReturnIPv4(v4, networkName, returnDate, notifyEMail string) (string, error) {
	// input check
	if v4 == "" {
		return "", fmt.Errorf("IPアドレスが指定されていません。")
	}
	if notifyEMail == "" {
		return "", fmt.Errorf("申請者メールアドレスが指定されていません。。")
	}
	if networkName == "" {
		return "", fmt.Errorf("ネットワーク名が指定されていません。。")
	}

	client, err := c.initAccess()
	if err != nil {
		return "", err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/assireturnv4regist.do?aplyid=108",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var actionURL string
	var token, destDisp, aplyId string

	// actionのURLを取得
	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "registconf") {
			actionURL = action
			doc.Find("input").Each(func(index int, s *goquery.Selection) {
				name, nameExists := s.Attr("name")
				value, valueExists := s.Attr("value")
				if nameExists && valueExists {
					switch name {
					case "org.apache.struts.taglib.html.TOKEN":
						token = value
					case "destdisp":
						destDisp = value
					case "aplyid":
						aplyId = value
					}
				}
			})
		}
	})

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	str := "org.apache.struts.taglib.html.TOKEN=" + token + "&destdisp=" + destDisp + "&aplyid=" + aplyId + "&ipaddr=" + v4 +
		"&netwrk_nm=" + networkName + "&rtn_date=" + returnDate +
		"&aply_from_addr=" + notifyEMail + "&aply_from_addr_confirm=" + notifyEMail + "&action=%90%5C%90%BF"
	// utf-8 => shift-jis
	reqBody, _, err := toShiftJIS(str)
	if err != nil {
		return "", err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// utf-8 => shift-jis
	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	// actionのURLを取得
	actionURL = ""
	token = ""
	var prevDispId string
	aplyId = ""
	destDisp = ""

	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "apply") {
			actionURL = action
			doc.Find("input").Each(func(index int, s *goquery.Selection) {
				name, nameExists := s.Attr("name")
				value, valueExists := s.Attr("value")
				if nameExists && valueExists {
					switch name {
					case "org.apache.struts.taglib.html.TOKEN":
						token = value
					case "prevDispId":
						prevDispId = value
					case "aplyid":
						aplyId = value
					case "destdisp":
						destDisp = value
					}
				}
			})
		}
	})

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	if strings.Contains(body, "IPネットワークアドレスが返却可能な割り当てアドレスではないか、ネットワーク名が正しくありません。") {
		return "", fmt.Errorf("IPネットワークアドレスが返却可能な割り当てアドレスではないか、ネットワーク名が正しくありません。")
	}

	if !strings.Contains(body, "上記の申請内容でよろしければ、「確認」ボタンを押してください。") {
		return "", fmt.Errorf("何かしらのエラーが発生しました。")
	}

	str = "org.apache.struts.taglib.html.TOKEN=" + token + "&prevDispId=" + prevDispId + "&aplyid=" + aplyId +
		"&destdisp=" + destDisp + "&inputconf=%8Am%94F"
	// utf-8 => shift-jis
	reqBody, _, err = toShiftJIS(str)
	if err != nil {
		return "", err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// utf-8 => shift-jis
	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var recepNo string

	// actionのURLを取得
	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						ok := false
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							if index == 0 && strings.Contains(tableCell2.Text(), "受付番号") {
								ok = true
							} else if index == 1 && ok {
								recepNo = tableCell2.Text()
							}
						})
					})
				})
			})
		})
	})

	return recepNo, nil
}

func (c *Config) ReturnIPv6(v6 []string, notifyEMail, returnDate string) (string, error) {
	// input check
	if len(v6) == 0 {
		return "", fmt.Errorf("IPアドレスが指定されていません。")
	}
	for _, ip := range v6 {
		if ip == "" {
			return "", fmt.Errorf("文字列が空のものがあります。")
		}
	}
	if notifyEMail == "" {
		return "", fmt.Errorf("申請者メールアドレスが指定されていません。。")
	}

	client, err := c.initAccess()
	if err != nil {
		return "", err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/G11220.do?aplyid=1106",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var actionURL string

	// actionのURLを取得
	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "Dispatch") {
			actionURL = action
		}
	})

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	//count := 0
	var returnIPv6List []ReturnIPv6List

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						var tmpIPv6List ReturnIPv6List
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							dataStr := strings.TrimSpace(tableCell2.Text())

							switch index {
							case 0:
								tmpIPv6List.NetworkID, _ = tableCell2.Find("input").Attr("value")
							case 1:
								tmpIPv6List.IPAddress = dataStr
							case 2:
								tmpIPv6List.NetworkName = dataStr
							case 3:
								tmpIPv6List.InfraUserKind = dataStr
							case 4:
								tmpIPv6List.AssignDate = dataStr
							}
						})
						returnIPv6List = append(returnIPv6List, tmpIPv6List)
					})
				})
			})
		})
	})

	var networkIDStr string

	for _, returnIPv6 := range returnIPv6List {
		for _, tmpIP := range v6 {
			if returnIPv6.IPAddress == tmpIP {
				if networkIDStr == "" {
					networkIDStr = "netwrkId=" + returnIPv6.NetworkID
				} else {
					networkIDStr += "&netwrkId=" + returnIPv6.NetworkID
				}
				break
			}
		}
	}

	if networkIDStr == "" {
		return "", fmt.Errorf("%s", "一致するNetworkIDがありません。")
	}

	str := "destdisp=G11220&aplyid=102&" + networkIDStr + "&action=%8Am%94F"
	// utf-8 => shift-jis
	reqBody, _, err := toShiftJIS(str)
	if err != nil {
		return "", err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	actionURL = ""

	// actionのURLを取得
	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "Dispatch") {
			actionURL = action
		}
	})

	str = "destdisp=G11221&aplyid=102&return_date=" +
		returnDate + "&aply_from_addr=" + notifyEMail + "&aply_from_addr_confirm=" + notifyEMail + "&action=%90%5C%90%BF"
	// utf-8 => shift-jis
	reqBody, _, err = toShiftJIS(str)
	if err != nil {
		return "", err
	}

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	if strings.Contains(body, "申請者メールアドレスを正しく入力してください") {
		return "", fmt.Errorf("JPNIC Response: 申請者メールアドレスを正しく入力してください")
	}

	if !strings.Contains(body, "上記の申請内容でよろしければ、｢確認｣ボタンを押してください。") {
		return "", fmt.Errorf("JPNIC Response: 何かしらのエラーが発生しています。")
	}

	// actionのURLを取得
	actionURL = ""

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "Dispatch") {
			actionURL = action
		}
	})

	str = "aplyid=102&inputconf=%8Am%94F"
	// utf-8 => shift-jis
	reqBody, _, err = toShiftJIS(str)
	if err != nil {
		return "", err
	}

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var recepNo string
	// utf-8 => shift-jis
	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	// actionのURLを取得
	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						ok := false
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							if index == 0 && strings.Contains(tableCell2.Text(), "受付番号") {
								ok = true
							} else if index == 1 && ok {
								recepNo = tableCell2.Text()
							}
						})
					})
				})
			})
		})
	})

	return recepNo, nil
}

func (c *Config) ChangeUserInfo(input JPNICHandleInput) (string, error) {
	client, err := c.initAccess()
	if err != nil {
		return "", err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/disppocregist.do?aplyid=501",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var actionURL string
	var token, destDisp, aplyId string

	// actionのURLを取得
	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "regist.do") {
			actionURL = action
			doc.Find("input").Each(func(index int, s *goquery.Selection) {
				name, nameExists := s.Attr("name")
				value, valueExists := s.Attr("value")
				if nameExists && valueExists {
					switch name {
					case "org.apache.struts.taglib.html.TOKEN":
						token = value
					case "destdisp":
						destDisp = value
					case "aplyid":
						aplyId = value
					}
				}
			})
		}
	})

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	// 初期値はJPNIC Handleで指定していた場合を想定
	kind := "person"
	if !input.IsJPNICHandle {
		// Group Handleで指定していた場合
		kind = "group"
	}

	str := "org.apache.struts.taglib.html.TOKEN=" + token + "&destdisp=" + destDisp + "&aplyid=" + aplyId +
		"&kind=" + kind + "&jpnic_hdl=" + input.JPNICHandle +
		"&name_jp=" + input.Name + "&name=" + input.NameEn + "&email=" + input.Email +
		"&org_nm_jp=" + input.Org + "&org_nm=" + input.OrgEn +
		"&zipcode=" + input.ZipCode + "&addr_jp=" + input.Address + "&addr=" + input.AddressEn +
		"&division_jp=" + input.Division + "&division=" + input.DivisionEn +
		"&title_jp=" + input.Title + "&title=" + input.TitleEn +
		"&phone=" + input.Tel + "&fax=" + input.Fax + "&ntfy_mail=" + input.NotifyMail +
		"&aply_from_addr=" + input.ApplyMail + "&aply_from_addr_confirm=" + input.ApplyMail + "&action=%90%5C%90%BF"

	// utf-8 => shift-jis
	reqBody, _, err := toShiftJIS(str)
	if err != nil {
		return "", err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// utf-8 => shift-jis
	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	// actionのURLを取得
	actionURL = ""
	token = ""
	var prevDispId string
	aplyId = ""
	destDisp = ""

	doc.Find("form").Each(func(_ int, formHtml *goquery.Selection) {
		action, _ := formHtml.Attr("action")
		if strings.Contains(action, "apply") {
			actionURL = action
			doc.Find("input").Each(func(index int, s *goquery.Selection) {
				name, nameExists := s.Attr("name")
				value, valueExists := s.Attr("value")
				if nameExists && valueExists {
					switch name {
					case "org.apache.struts.taglib.html.TOKEN":
						token = value
					case "prevDispId":
						prevDispId = value
					case "aplyid":
						aplyId = value
					case "destdisp":
						destDisp = value
					}
				}
			})
		}
	})

	if actionURL == "" {
		return "", fmt.Errorf("action URLの取得失敗")
	}

	if strings.Contains(body, "IPネットワークアドレスが返却可能な割り当てアドレスではないか、ネットワーク名が正しくありません。") {
		return "", fmt.Errorf("IPネットワークアドレスが返却可能な割り当てアドレスではないか、ネットワーク名が正しくありません。")
	}

	if strings.Contains(body, "JPNICハンドルを正しく入力してください") {
		return "", fmt.Errorf("正しいJPNICハンドルを入力してください")
	}

	if !strings.Contains(body, "上記の申請内容でよろしければ、「確認」ボタンを押してください。") {
		return "", fmt.Errorf("何かしらのエラーが発生しました。")
	}

	str = "org.apache.struts.taglib.html.TOKEN=" + token + "&prevDispId=" + prevDispId + "&aplyid=" + aplyId +
		"&destdisp=" + destDisp + "&inputconf=%8Am%94F"
	// utf-8 => shift-jis
	reqBody, _, err = toShiftJIS(str)
	if err != nil {
		return "", err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp" + actionURL,
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// utf-8 => shift-jis
	body, _, err = readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var recepNo string

	// actionのURLを取得
	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(_ int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						ok := false
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							if index == 0 && strings.Contains(tableCell2.Text(), "受付番号") {
								ok = true
							} else if index == 1 && ok {
								recepNo = tableCell2.Text()
							}
						})
					})
				})
			})
		})
	})

	return recepNo, nil
}

func (c *Config) GetRequestList(searchStr string) ([]RequestInfo, error) {
	client, err := c.initAccess()
	if err != nil {
		return nil, err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	str := "destdisp=D10009&startRecepNo=" + searchStr + "&endRecepNo=&deliNo=&aplyKind=&aplyClass=&resceAdmSnm=&aplyDateS=&aplyDateE=&completDateS=&completDateE=&statusId=&pswdResceNewConfirm=%81%40%8C%9F%8D%F5%81%40"
	// utf-8 => shift-jis
	reqBody, _, err := toShiftJIS(str)
	if err != nil {
		return nil, err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/applysearchlink.do",
		Body:        reqBody,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.post()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	count := 0
	var infos []RequestInfo

	doc.Find("table").Each(func(_ int, tableHtml *goquery.Selection) {
		tableHtml.Find("tr").Each(func(_ int, rowHtml *goquery.Selection) {
			var info RequestInfo
			rowHtml.Find("td").Each(func(index int, tableCell *goquery.Selection) {
				dataStr := strings.TrimSpace(tableCell.Text())

				switch index {
				case 0:
					info.RecepNo = dataStr
				case 1:
					info.DeliNo = dataStr
				case 2:
					info.ApplyKind = dataStr
				case 3:
					info.ApplyClass = dataStr
				case 4:
					info.Applicant = dataStr
				case 5:
					info.ApplyDate = dataStr
				case 6:
					info.CompleteDate = dataStr
				case 7:
					info.Status = dataStr
				}
				count++
			})
			if count == 8 {
				infos = append(infos, info)
			}
			count = 0
		})
	})

	infos = infos[1:]

	return infos, nil
}

func (c *Config) GetDetailRequest(recepNo string) (string, error) {
	client, err := c.initAccess()
	if err != nil {
		return "", err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/applyform.do?recepNo=" + recepNo,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var info string

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(index int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							dataStr := strings.TrimSpace(tableCell2.Text())
							if dataStr != "" {
								info = "\n" + dataStr
							}
						})
					})
				})
			})
		})
	})

	return info, nil
}

func (c *Config) GetResourceManagement() (ResourceInfo, error) {
	var info ResourceInfo

	client, err := c.initAccess()
	if err != nil {
		return info, err
	}

	r := request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	url, err := getLink(client, "資源管理者情報")
	if err != nil {
		return info, err
	}

	r = request{
		Client:      client,
		URL:         "https://iphostmaster.nic.ad.jp/jpnic/" + url,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err = r.get()
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	body, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return info, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return info, err
	}

	re := regexp.MustCompile(`\(([^}]*)\)`)
	err = nil

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(index int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							tableCell2.Find("table").Each(func(_ int, tableHtml3 *goquery.Selection) {
								tableHtml3.Find("tr").Each(func(_ int, rowHtml3 *goquery.Selection) {
									rowHtml3.Find("td").Each(func(index int, tableCell3 *goquery.Selection) {
										tableCell3.Find("table").Each(func(_ int, tableHtml4 *goquery.Selection) {
											tableHtml4.Find("tr").Each(func(_ int, rowHtml4 *goquery.Selection) {
												var title string
												cidrBlockSegment := false
												var cidrBlock ResourceCIDRBlock
												rowHtml4.Find("td").Each(func(index int, tableCell4 *goquery.Selection) {
													dataStr := strings.TrimSpace(tableCell4.Text())

													switch index {
													case 0:
														cidrBlockSegment = false
														title = dataStr
														addressDetailURL, addressExists := tableCell4.Find("a").Attr("href")
														if addressExists {
															cidrBlockSegment = strings.Contains(addressDetailURL, "entryinfo")
															splitAddress := strings.Split(dataStr, "(")
															tmpAddress := strings.Replace(splitAddress[0], "\n", "", 1)
															address := strings.Replace(tmpAddress, "	", "", 3)
															cidrBlock.Address = strings.TrimSpace(address)
														}
													case 1:
														switch title {
														case "資源管理者番号":
															info.ResourceManagerInfo.ResourceManagerNo = dataStr
														case "資源管理者略称":
															info.ResourceManagerInfo.Ryakusyo = dataStr
														case "管理組織名":
															info.ResourceManagerInfo.Org = dataStr
														case "Organization":
															info.ResourceManagerInfo.OrgEn = dataStr
														case "郵便番号":
															info.ResourceManagerInfo.ZipCode = dataStr
														case "住所":
															info.ResourceManagerInfo.Address = dataStr
														case "Address":
															info.ResourceManagerInfo.AddressEn = dataStr
														case "電話番号":
															info.ResourceManagerInfo.Tel = dataStr
														case "FAX番号":
															info.ResourceManagerInfo.Fax = dataStr
														case "資源管理責任者":
															info.ResourceManagerInfo.ResourceManagementManager = dataStr
														case "連絡担当窓口":
															info.ResourceManagerInfo.ContactPerson = dataStr
														case "一般問い合わせ窓口":
															info.ResourceManagerInfo.Inquiry = dataStr
														case "資源管理者通知アドレス":
															info.ResourceManagerInfo.NotifyMail = dataStr
														case "アサインメントウィンドウサイズ":
															info.ResourceManagerInfo.AssigmentWindowSize = dataStr
														case "管理開始日":
															info.ResourceManagerInfo.ManagementStartDate = dataStr
														case "管理終了日":
															info.ResourceManagerInfo.ManagementEndDate = dataStr
														case "最終更新日":
															info.ResourceManagerInfo.UpdateDate = dataStr
														default:
															if cidrBlockSegment {
																cidrBlock.AssignDate = dataStr
															}
														}
													case 2:
														switch title {
														case "総利用率":
															match := re.FindStringSubmatch(dataStr)
															if len(match) == 0 {
																err = fmt.Errorf("データが存在しません")
																break
															}
															splitAddress := strings.Split(match[1], "/")

															info.UsedAddress, err = strconv.ParseUint(splitAddress[0], 10, 32)
															if err != nil {
																break
															}
															info.AllAddress, err = strconv.ParseUint(splitAddress[1], 10, 32)
															if err != nil {
																break
															}

															info.UtilizationRatio, err = strconv.ParseFloat(dataStr[:strings.Index(dataStr, "%")], 16)
															if err != nil {
																break
															}
														case "ＡＤ　ｒａｔｉｏ":
															log.Println(strconv.Itoa(index) + ": " + dataStr)

															info.ADRatio, err = strconv.ParseFloat(dataStr, 16)
															if err != nil {
																break
															}
														default:
															if cidrBlockSegment {
																match := re.FindStringSubmatch(dataStr)
																if len(match) == 0 {
																	err = fmt.Errorf("データが存在しません")
																	break
																}
																splitAddress := strings.Split(match[1], "/")

																cidrBlock.UsedAddress, err = strconv.ParseUint(splitAddress[0], 10, 32)
																if err != nil {
																	break
																}
																cidrBlock.AllAddress, err = strconv.ParseUint(splitAddress[1], 10, 32)
																if err != nil {
																	break
																}

																cidrBlock.UtilizationRatio, err = strconv.ParseFloat(dataStr[:strings.Index(dataStr, "%")], 16)
																if err != nil {
																	break
																}
															}
														}
													}
													if cidrBlockSegment && index == 2 {
														info.ResourceCIDRBlock = append(info.ResourceCIDRBlock, cidrBlock)
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

	if err != nil {
		return info, err
	}
	return info, nil
}
