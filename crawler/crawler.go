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
var noSucReq = 0

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
	fmt.Println("created client proxy successfully")
	crawler := &Crawler{
		Client: client,
		Delay:  1 * time.Second,
	}

	err = crawler.Crawl()
	if err != nil {
		return err
	}
	fmt.Println("done crawling")
	return nil
}

func (crwl *Crawler) Crawl() error {
	pages := crwl.FindLinks(BaseURL)
	for _, page := range pages {

		//swapping IP
		newClient, _ := getProxyClient()
		crwl.Client = newClient

		crwl.GetFragrances(page) // need to implement concurrency in future
		time.Sleep(crwl.Delay)
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
	fmt.Println("found all links in base URL")
	return resultArr
}

func (crwl *Crawler) GetFragrances(url string) {
	var err error
	res, err := CreateRequest(crwl.Client, BaseBaseURL+url)
	if err != nil {
		log.Fatal(err) //429 too many requests
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return
	}
	noSucReq++ //to delete in future

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	name := doc.Find("h1").Text()
	brand := doc.Find("span[itemprop='name'][class='vote-button-name']").Text()
	fmt.Println(res.StatusCode, name, brand)
	//mainAccords := doc.Find("h6").Text()
	//
	//doc.Find("h4[style='margin-top: 0.5rem;']").Each(func(i int, s *goquery.Selection) {
	//	s.Next().Children().Each(func(i int, s *goquery.Selection) {
	//		//noteName := s.Children().Last().Text()
	//	}) //this div contains all divs containing notes

	//})
	//
	//doc.Find(".accord-bar").Each(func(i int, s *goquery.Selection) {
	//	fmt.Println("accord: ", s.Text())
	//})
	////notes := doc.Find("div#pyramid")
	//fmt.Println(name, brand, mainAccords)
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
