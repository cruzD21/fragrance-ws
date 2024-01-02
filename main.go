package main

import "fragrance-ws/crawler"

var BaseURL = "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"

func main() {

	crwl := &crawler.Crawler{}

	perfumeLinks := crwl.FindLinks(BaseURL)
	for _, link := range perfumeLinks {
		crwl.GetFragrances(link)
	}
}
