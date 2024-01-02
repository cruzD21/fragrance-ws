package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

var BaseBaseURL = "https://www.fragrantica.com"

type Crawler struct{}

func (crwl *Crawler) FindLinks(baseURL string) []string {
	var resultArr []string

	client := CreateClient(nil)
	res, err := CreateRequest(client, baseURL)
	if err != nil {
		log.Fatal(err)
	}

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		//perfumeName := item.Text()
		//could make this a function where prevents /perfume-review
		if validURL(link) {
			resultArr = append(resultArr, link)
			//fmt.Println("Perfume Name: ", perfumeName)
			//fmt.Println("url: ", url)

		}
	})
	fmt.Println("got all links")
	return resultArr
}

func (crwl *Crawler) GetFragrances(url string) {
	var err error
	client := CreateClient(nil)
	res, err := CreateRequest(client, BaseBaseURL+url)

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
