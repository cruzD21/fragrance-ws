package crawler

import (
	"errors"
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
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("Error Sending Request, Status Code: %d , %s", res.StatusCode, res.Status)
		return nil, errors.New(errStr)
	}
	return res, nil
}

func randomUserAgent() string {
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}
