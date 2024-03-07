package main

import (
	"fragrance-ws/crawler"
	"log"
)

func main() {

	if err := crawler.Run(); err != nil {
		log.Fatal(err)
	}

	//crawler.Test()
}
