package crawler

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
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

func TestProxy() {

	for i := 0; i < 5; i++ {
		client, _ := getProxyClient()
		res, err := CreateRequest(client, "https://lumtest.com/myip.json")
		if err != nil {
			log.Fatal(err)
		}
		//printing ip logic

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Printing the response body
		fmt.Println("Status Code:", res.StatusCode)
		fmt.Println("Response Body:", string(body))

		// You can also print other parts of the response like status code, headers etc.

		//fmt.Println("Headers:", res.Header)
		res.Body.Close()
	}
}
