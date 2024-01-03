package crawler

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
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
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	return createClient(parsedURL), nil
}
