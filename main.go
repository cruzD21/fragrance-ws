package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var BaseURL = "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"
var BaseBaseURL = "https://www.fragrantica.com"
var maxRequests = 7
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

type Collector struct {
	reqNo       int
	Delay       time.Duration
	RandomDelay time.Duration
	SleepTime   time.Duration
}

func CreateClient(proxyString interface{}) *http.Client {
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
					DualStack: true,
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
					DualStack: true,
				}).DialContext,
			},
		}
	}
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

//func createClient() (*http.Client, error) {
//	// creates clients and sets up transport
//	var err error
//	client := &http.Client{
//		Timeout: 15 * time.Second,
//		Transport: &http.Transport{
//			DialContext: (&net.Dialer{
//				Timeout:   15 * time.Second,
//				KeepAlive: 15 * time.Second,
//				DualStack: true,
//			}).DialContext,
//		},
//	}
//	// Set the client Transport to the RoundTripper that solves the Cloudflare anti-bot
//	client.Transport, err = cfrt.New(client.Transport)
//	return client, err
//
//}

func validURL(url string) bool {
	if !strings.HasPrefix(url, "/perfume") {
		return false
	}
	if strings.HasPrefix(url, "/perfume-review") || strings.HasPrefix(url, "/perfume-finder") {
		return false
	}
	//url only has strictly "/perfume"
	return true
}

func findLinks(baseURL string) []string {
	var links []string
	client := CreateClient(nil)
	req, _ := http.NewRequest("GET", baseURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making req find links", err)
	}
	defer res.Body.Close()
	fmt.Println("res tatus from findLinks", res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		//perfumeName := item.Text()
		//could make this a function where prevents /perfume-review
		if validURL(link) {
			links = append(links, link)
			//fmt.Println("Perfume Name: ", perfumeName)
			//fmt.Println("url: ", url)

		}
	})
	fmt.Println("got all links")
	return links
}

func getFragrances(url string) {
	var err error
	client := CreateClient(nil)

	req, err := http.NewRequest("GET", BaseBaseURL+url, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	if err != nil {
		log.Println("error creating request", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("statos", res.StatusCode)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Status Code Error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
	h6 := doc.Find("h6").Text()
	fmt.Println(h6)
	doc.Find(".accord-bar").Each(func(i int, s *goquery.Selection) {
		fmt.Println("accord: ", s.Text())
	})

	fmt.Println("sleeping 10 sec before next req")
	time.Sleep(10 * time.Second)
}

func main() {

	perfumeLinks := findLinks(BaseURL)
	for _, link := range perfumeLinks {
		getFragrances(link)
	}

}
