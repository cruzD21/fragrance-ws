package crawler

import (
	"fmt"
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

	client, _ := getProxyClient()
	res, err := CreateRequest(client, BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	//printing ip logic

	frag, err := parseFragrancePage(res)

	// Printing the response body

	// You can also print other parts of the response like status code, headers etc.

	//fmt.Println("Headers:", res.Header)
	res.Body.Close()

}
