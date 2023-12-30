package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var operatingSystem = []string{
	"Windows", "macOS",
}

var userAgents = map[string][]string{
	"Windows": {
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	},
	"macOS": {
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	},
}

var visited = make(map[string]bool)

func randomOS() (string, string) {
	index := rand.Intn(len(operatingSystem))
	osIndex := rand.Intn(3)
	ops := operatingSystem[index]

	return ops, userAgents[ops][osIndex]
}

func ChangeHeaders(h *http.Header) {
	p, userAgent := randomOS()
	h.Set("sec-ch-ua-platform", p)
	h.Set("user-agent", userAgent)
}

var nDelay = 1 * time.Second
var rDelay = 2 * time.Second
var waitTime = 1 * time.Minute

func main() {
	fakeChrome := req.DefaultClient().ImpersonateChrome()
	baseURL := "https://www.fragrantica.com/perfume/Le-Labo/Another-13-10131.html"

	maxReq := 9
	noFrag := 0
	reqCount := 0
	//res, e := fakeChrome.R().Get(baseURL)
	//if e != nil {
	//	log.Fatal(e)
	//}
	//fmt.Println(res.StatusCode)

	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	)
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       nDelay,
		RandomDelay: rDelay,
	})
	if err != nil {
		log.Fatal(err)
	}

	c.SetClient(&http.Client{
		Transport: fakeChrome.Transport,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Page: ", r.URL)
		if visited[r.URL.String()] {
			fmt.Println("already visited page", r.URL.String())
			r.Abort()
		}
		reqCount++
		if reqCount >= maxReq {
			fmt.Printf("max amount of req reached, sleeping for %s min . . .\n", waitTime)
			time.Sleep(waitTime)
			fmt.Println("resuming requests")
			waitTime += 1 * time.Minute
			reqCount = 0
		}
		//change headers logic
		//fmt.Println("previous headers", r.Headers)
		ChangeHeaders(r.Headers)
		//fmt.Println("new headers", r.Headers)

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Status Code: ", r.StatusCode)
		fmt.Println("request  headers ", r.Request.Headers)
		fmt.Println("response headers ", r.Headers)
		log.Fatal(err)

	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Status Code:", r.StatusCode)
		visited[r.Request.URL.String()] = true
	})
	//
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/perfume") {
			absoluteURL := e.Request.AbsoluteURL(link)
			fmt.Println("Found perfume link:", absoluteURL)
			e.Request.Visit(absoluteURL)
		}
	})

	c.OnHTML("div#main-content", func(e *colly.HTMLElement) {
		e.DOM.Find("h1").Each(func(i int, s *goquery.Selection) {
			fmt.Println(s.Text())
		})
		h6 := e.DOM.Find("h6").Text()
		fmt.Println(h6)
		e.DOM.Find(".accord-bar").Each(func(i int, s *goquery.Selection) {
			fmt.Println("accord: ", s.Text())
		})
		noFrag++
		fmt.Println("fragNo", noFrag)
	})

	er := c.Visit(baseURL)
	if er != nil {
		log.Fatal(er)
	}

}
