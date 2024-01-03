package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

var BaseURL = "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"
var BaseBaseURL = "https://www.fragrantica.com"

type Crawler struct {
	Client *http.Client
	Delay  time.Duration
}

func Run() error {
	var err error
	//here I have a client ready to start crawling
	client, err := getProxyClient()
	if err != nil {
		return err
	}
	crawler := &Crawler{
		Client: client,
		Delay:  1 * time.Second,
	}

	err = crawler.Crawl()
	if err != nil {
		return err
	}

	return nil
}

func (crwl *Crawler) Crawl() error {
	pages := crwl.FindLinks(BaseURL)
	for _, page := range pages {
		crwl.GetFragrances(page)
		time.Sleep(crwl.Delay) // need to implement concurrency in future
	}

	return nil
}

func (crwl *Crawler) FindLinks(baseURL string) []string {
	var resultArr []string

	res, err := CreateRequest(crwl.Client, baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		if validURL(link) {
			resultArr = append(resultArr, link)
		}
	})
	return resultArr
}

func (crwl *Crawler) GetFragrances(url string) {
	var err error
	res, err := CreateRequest(crwl.Client, BaseBaseURL+url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

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
}

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
