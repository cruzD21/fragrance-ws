package crawler

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func createClient(proxyString interface{}) *http.Client {
	switch v := proxyString.(type) {

	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
				DialContext: (&net.Dialer{
					Timeout:   15 * time.Second,
					KeepAlive: 15 * time.Second,
					//DualStack: true,
				}).DialContext,
			},
		}
	default:
		return &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   15 * time.Second,
					KeepAlive: 15 * time.Second,
					//DualStack: true,
				}).DialContext,
			},
		}
	}
}

func CreateRequest(c *http.Client, url string) (*http.Response, error) {
	var err error

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(req)

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusForbidden {
		return res, nil //invalid ip, skip to next
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error Sending Request, Status Code: %d , %s", res.StatusCode, res.Status)
	}
	return res, nil
}

func randomUserAgent() string {
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func setHeaders(req *http.Request) {
	headers := http.Header{
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
		"Accept-Language":           {"en"},
		"Sec-Ch-Ua":                 {"\"Not A Brand\";v=\"99\", \"Chromium\";v=\"90\""},
		"Sec-Ch-Ua-Mobile":          {"?0"},
		"Sec-Ch-Ua-Platform":        {"\"Linux\""},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"},
	}
	req.Header = headers
}
