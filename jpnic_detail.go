package jpnic

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func getInfoDetail(client *http.Client, userURL string) (InfoDetail, error) {
	var info InfoDetail

	r := request{
		Client:      client,
		URL:         baseURL + userURL,
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		log.Println(err)
		return info, err
	}

	respBody, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return info, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(respBody))
	if err != nil {
		return info, err
	}

	var title string
	isTitle := true

	doc.Find("table").Children().Find("table").Children().Find("table").Children().Find("table").Children().Find("td").Each(func(_ int, tableHtml1 *goquery.Selection) {
		dataStr := strings.TrimSpace(tableHtml1.Text())
		if isTitle {
			title = dataStr
		}

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
			info.AdminJPNICHandleLink, _ = tableHtml1.Find("a").Attr("href")
		case "技術連絡担当者":
			info.TechJPNICHandle = dataStr
			info.TechJPNICHandleLink, _ = tableHtml1.Find("a").Attr("href")
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

		isTitle = !isTitle
	})

	return info, err
}

func getJPNICHandle(client *http.Client, handleURL string) (JPNICHandleDetail, error) {
	var info JPNICHandleDetail

	r := request{
		Client:      client,
		URL:         baseURL + "/jpnic/" + handleURL,
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return info, err
	}

	resBody, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return info, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resBody))
	if err != nil {
		return info, err
	}

	var title string
	isTitle := true

	doc.Find("table").Children().Find("table").Children().Find("table").Children().Find("td").Each(func(_ int, tableHtml1 *goquery.Selection) {
		dataStr := strings.TrimSpace(tableHtml1.Text())
		if isTitle {
			title = dataStr
		}

		switch title {
		case "グループハンドル":
			info.IsJPNICHandle = false
			info.JPNICHandle = dataStr
		case "グループ名":
			info.Name = dataStr
		case "Group Name":
			info.NameEn = dataStr
		case "JPNICハンドル":
			info.IsJPNICHandle = true
			info.JPNICHandle = dataStr
		case "氏名":
			info.Name = dataStr
		case "Last, First":
			info.NameEn = dataStr
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

		isTitle = !isTitle
	})

	return info, err
}

func getRecepDetail(client *http.Client, recepURL string) (string, error) {
	r := request{
		Client:      client,
		URL:         baseURL + recepURL,
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
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

	doc.Find("table").Children().Find("table").Children().Find("td").Each(func(index int, tableHtml1 *goquery.Selection) {
		dataStr := strings.TrimSpace(tableHtml1.Text())
		if dataStr != "" {
			info = "\n" + dataStr
		}
	})

	return info, nil
}
