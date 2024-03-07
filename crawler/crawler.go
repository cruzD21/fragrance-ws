package crawler

import (
	"fmt"
	"fragrance-ws/db"
	"fragrance-ws/models"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	DBConn = &db.DatabaseConn{}
)

var BaseURL = "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"
var BaseBaseURL = "https://www.fragrantica.com"
var noReq = 0

type Crawler struct {
	Client *http.Client
	DB     *db.DatabaseConn
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

	if err = crawler.Crawl(); err != nil {
		return err
	}

	return nil
}

func (crwl *Crawler) Crawl() error {
	var err error
	pages := crwl.FindLinks(BaseBaseURL)
	fmt.Println("printing pages")
	for i, page := range pages {
		fmt.Printf("page no %d, url: %s\n", i, page)
	}
	//fmt.Println(pages)
	err = DBConn.DatabaseInit()
	if err != nil {
		return err
	}

	for _, page := range pages {
		if noReq == 20 {
			return fmt.Errorf("max req done")
		}
		//swapping IP
		newClient, _ := getProxyClient()
		crwl.Client = newClient

		data, err := crwl.GetFragrances(page) // need to implement concurrency in future
		if err != nil {
			return err
		}

		err = DBConn.InsertPage(data)
		if err != nil {
			return err
		}
		noReq++
		//time.Sleep(crwl.Delay)
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

func (crwl *Crawler) GetFragrances(url string) (models.FragrancePage, error) {
	var err error
	res, err := CreateRequest(crwl.Client, BaseBaseURL+url)
	if err != nil {
		log.Fatal(err) //429 too many requests
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return models.FragrancePage{}, nil
	}

	data, err := parseFragrancePage(res)
	if err != nil {
		return models.FragrancePage{}, err
	}

	// add to db
	return data, nil
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
