package crawler

import (
	"errors"
	"fragrance-ws/models"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type Fragrance struct {
	Name        string
	Brand       string
	MainAccords []string
	TopNotes    []string
	MiddleNotes []string
	BaseNotes   []string
}

func parseFragrancePage(res *http.Response) (models.FragrancePage, error) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("problem in goquery")
		return models.FragrancePage{}, err
	}

	fragrance, err := parseFragrance(doc)
	if err != nil {
		log.Println("problem in frag")
		return models.FragrancePage{}, err
	}
	fragranceHouse, err := parseFragranceHouse(doc)
	if err != nil {
		log.Println("error in house")
		return models.FragrancePage{}, err
	}
	noteCategories := parseNotePyramid(doc)
	accords := parseAccords(doc)

	return models.FragrancePage{
		Fragrance:   fragrance,
		FragHouse:   fragranceHouse,
		NoteCat:     noteCategories,
		MainAccords: accords,
	}, nil
}

func parseFragrance(doc *goquery.Document) (models.Fragrance, error) {
	name := doc.Find("h1").Text()
	description := doc.Find("div[itemprop='description'] p").First().Text()
	log.Println(name)

	if name == "" {
		log.Println("empty frag name im skipping ")
		return models.Fragrance{}, errors.New("the fragrance name is empty")
	}

	if description == "" {
		return models.Fragrance{}, errors.New("the fragrance description is empty")
	}

	return models.Fragrance{
		Name:        name,
		Description: description,
	}, nil
}

func parseFragranceHouse(doc *goquery.Document) (models.FragranceHouse, error) {
	name := doc.Find(".vote-button-name").Text()
	log.Println(name)
	if len(name) == 0 {
		log.Println("empty house name im skipping ", name)
		return models.FragranceHouse{}, errors.New("empty fragrance name")
	}
	return models.FragranceHouse{
		Name: name,
	}, nil
}

func parseNotePyramid(doc *goquery.Document) models.NoteCategories {
	//some notes only have one level so need to account for that
	var topNotes, middleNotes, baseNotes []string

	doc.Find("h4 > b").Each(func(i int, s *goquery.Selection) {
		pyramidLevel := s.Text() //pyramid accords
		s.Parent().Next().Contents().Children().Each(func(i int, s2 *goquery.Selection) {
			note := s2.Children().Last().Last().Text()
			switch pyramidLevel {
			case "Top Notes":
				topNotes = append(topNotes, note)
			case "Middle Notes":
				middleNotes = append(middleNotes, note)
			case "Base Notes":
				baseNotes = append(baseNotes, note)
			}
		})
	})
	return models.NoteCategories{
		TopNotes:    topNotes,
		MiddleNotes: middleNotes,
		BaseNotes:   baseNotes,
	}
}

func parseAccords(doc *goquery.Document) []string {
	var mainAccords []string
	doc.Find(".accord-bar").Each(func(i int, s *goquery.Selection) {
		accord := s.Text()
		mainAccords = append(mainAccords, accord)
	})
	return mainAccords
}
