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
	res, err := CreateRequest(client, "https://www.fragrantica.com/perfume/Puzzle-Parfum/Puzzle-Daylight-91072.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	//printing ip logic

	page, err := parseFragrancePage(res)
	if err != nil {
		log.Printf("error inserting into db with error : %e", err)
	}
	fmt.Printf("%+v\n", page)

	err = supa.InsertPage(page)
	if err != nil {
		log.Printf("error inserting into db with error : %e", err)
	}

	log.Println("code inserted  page successfully to db")
}
