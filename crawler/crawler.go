package crawler

import (
	"fmt"
	"fragrance-ws/db"
	"fragrance-ws/models"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"path"
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
	var failedUrls []string
	var data models.FragrancePage
	var fragID string

	pages := crwl.FindLinks(BaseBaseURL) // no repeating urls saved
	log.Println("total pages 	", len(pages))

	//initialize external services
	err = DBConn.DatabaseInit()
	if err != nil {
		return err
	}
	awsClient, err := awsInit()
	if err != nil {
		log.Printf("awsInit -> %v", err)
		return err
	}

	for page := range pages {
		//swapping IP
		newClient, _ := getProxyClient()
		crwl.Client = newClient

		data, fragID, err = crwl.GetFragrances(page)
		if err != nil {
			log.Printf("crwl.GetFragances -> %v, skipping page", err)
			continue
		}
		fmt.Printf("%+v\n", data)

		err = DBConn.InsertPage(data)
		if err != nil {
			log.Printf("DBConn.InsertPage -> %v", err)
			failedUrls = append(failedUrls, page)
			continue //skip to next page
		}
		err = putObjectIntoBucket(awsClient, fragID) //insert fragrance image
		if err != nil {
			log.Printf("putObjectIntoBucket ->  %v, image not inserted", err)
			continue
		}

	}

	return nil
}

func (crwl *Crawler) FindLinks(baseURL string) map[string]bool {
	resultMap := make(map[string]bool)

	res, err := CreateRequest(crwl.Client, baseURL)
	if err != nil {

		log.Fatal(err) //
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		if validURL(link) {
			if _, exist := resultMap[link]; !exist {
				resultMap[link] = true
			}
		}
	})
	fmt.Println("found all links in base URL")
	return resultMap
}

func (crwl *Crawler) GetFragrances(url string) (models.FragrancePage, string, error) {
	var err error
	var data models.FragrancePage

	res, err := CreateRequest(crwl.Client, BaseBaseURL+url)
	if err != nil {
		log.Fatal(err) //429 too many requests
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return models.FragrancePage{}, "", nil
	}

	data, err = parseFragrancePage(res)
	if err != nil {
		log.Printf("parseFragrancePage -> %v", err)
		return models.FragrancePage{}, "", err
	}

	fragID, err := extractID(url)
	if err != nil {
		log.Printf("extractID -> %v", err)
		return models.FragrancePage{}, "", err
	}
	data.Fragrance.FraganticaID = fragID

	// add to db
	return data, fragID, nil
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

func extractID(url string) (string, error) {
	lastPart := path.Base(url)

	lastDotIndex := strings.LastIndex(lastPart, ".")
	if lastDotIndex == -1 {
		return "", fmt.Errorf("no extension found in URL")
	}

	namePart := lastPart[:lastDotIndex]
	lastNonDigitIndex := strings.LastIndexFunc(namePart, func(r rune) bool {
		return !('0' <= r && r <= '9')
	})

	digits := namePart[lastNonDigitIndex+1:]
	return digits, nil
}
