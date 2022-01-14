package jpnic

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type request struct {
	Client          *http.Client
	URL             string
	Body            string
	UserAgent       string
	ContentType     string
	SessionID       string
	ServerSessionID string
}

func (c *Config) initAccess(menuName string) (*http.Client, string, error) {
	sessionID, err := randomStr()
	if err != nil {
		return nil, "", err
	}

	cookies := []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: sessionID,
		},
	}

	// Cookie
	urlObj, _ := url.Parse("https://iphostmaster.nic.ad.jp/")
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, "", err
	}

	jar.SetCookies(urlObj, cookies)

	// Load .p12 File
	p12Bytes, err := ioutil.ReadFile(c.PfxFilePath)
	if err != nil {
		return nil, "", err
	}

	// .p12 decode
	key, cert, err := pkcs12.Decode(p12Bytes, c.PfxPass)
	if err != nil {
		return nil, "", err
	}

	// Load CA
	caCertBytes, err := ioutil.ReadFile(c.CAFilePath)
	if err != nil {
		return nil, "", err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key,
			Leaf:        cert,
		}},
		RootCAs: caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Jar: jar}

	// Login
	r := request{
		Client:      client,
		URL:         baseURL + "/jpnic/certmemberlogin.do",
		Body:        "",
		UserAgent:   userAgent,
		ContentType: contentType,
	}

	resp, err := r.get()
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	result, _, err := readShiftJIS(resp.Body)
	if err != nil {
		return nil, "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result))
	if err != nil {
		return nil, "", err
	}
	resultContent, isExists := doc.Find("meta").Attr("content")
	if !isExists {
		return nil, "", fmt.Errorf("エラーが発生しました")
	}
	refreshURL := strings.Split(resultContent, "=")[1]

	// menu
	menuURL, err := getLink(client, refreshURL, menuName)
	if err != nil {
		return nil, "", err
	}

	return client, menuURL, nil
}

func getJSessionID(cookies []*http.Cookie) string {
	var jsessionID string

	for _, tmp := range cookies {
		if tmp.Name == "jsessionid" {
			jsessionID = tmp.Value
			break
		}
		//log.Println(tmp.Name)
		//log.Println(tmp.Value)
	}
	return jsessionID
}

func (r *request) post() (*http.Response, error) {
	req, err := http.NewRequest("POST", r.URL, ioutil.NopCloser(bytes.NewBufferString(r.Body)))

	req.Header.Add("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Host", "iphostmaster.nic.ad.jp")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Upgrade-Insecure-Requests","1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")

	resp, err := r.Client.Do(req)

	return resp, err
}

func (r *request) get() (*http.Response, error) {
	req, err := http.NewRequest("GET", r.URL, nil)

	req.Header.Add("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Host", "iphostmaster.nic.ad.jp")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Upgrade-Insecure-Requests","1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")

	resp, err := r.Client.Do(req)
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 503:
		body, _, _ := readEUCJP(resp.Body)
		// メンテナンス判定
		if strings.Contains(body, "ただいまメンテナンス中です") {
			return resp, fmt.Errorf("[%d] 現在、メンテナンス中のためデータ取得が出来ません。", resp.StatusCode)
		} else {
			return resp, fmt.Errorf("Status Code: %d ", resp.StatusCode)
		}
	}

	return resp, err
}
