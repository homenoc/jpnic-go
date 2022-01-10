package jpnic

import (
	"crypto/rand"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
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

func readEUCJP(readerInput io.Reader) (string, []byte, error) {
	// euc-jp => utf-8
	reader := transform.NewReader(readerInput, japanese.EUCJP.NewDecoder())
	strByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", nil, err

	}
	return string(strByte), strByte, nil
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

func getLink(client *http.Client, menuURL, str string) (string, error) {
	r := request{
		Client:      client,
		URL:         baseURL + "/jpnic/" + menuURL,
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

	var url string

	doc.Find("table").Each(func(_ int, tableHtml1 *goquery.Selection) {
		tableHtml1.Find("tr").Each(func(_ int, rowHtml1 *goquery.Selection) {
			rowHtml1.Find("td").Each(func(index int, tableCell1 *goquery.Selection) {
				tableCell1.Find("table").Each(func(_ int, tableHtml2 *goquery.Selection) {
					tableHtml2.Find("tr").Each(func(_ int, rowHtml2 *goquery.Selection) {
						rowHtml2.Find("td").Each(func(index int, tableCell2 *goquery.Selection) {
							tableCell2.Find("a").Each(func(index int, aCell1 *goquery.Selection) {
								dataStr := strings.TrimSpace(aCell1.Text())
								if strings.Contains(dataStr, str) {
									url, _ = aCell1.Attr("href")
								}
							})
						})
					})
				})
			})
		})
	})

	if url == "" {
		return "", fmt.Errorf("項目が見つかりません")
	}

	return url, nil
}

func getSearchBoolean(isFilter bool) string {
	if isFilter {
		return "on"
	} else {
		return "off"
	}
}
