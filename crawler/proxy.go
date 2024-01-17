package crawler

import (
	"fmt"
	"fragrance-ws/db"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func getProxyClient() (*http.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	KEY := os.Getenv("KEY")
	HOST := os.Getenv("HOST")
	proxyURL := fmt.Sprintf("http://%s@%s", KEY, HOST)
	return createClient(proxyURL), nil
}

func Test() {
	supa := db.DatabaseConn{}
	err := supa.DatabaseInit()

	client, _ := getProxyClient()
	res, err := CreateRequest(client, BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	//printing ip logic

	page, err := parseFragrancePage(res)

	err = supa.InsertPage(page)
	if err != nil {
		log.Fatalf("error inserting into db with error : %e", err)
	}

	log.Println("code inserted  page successfully to db")
}
