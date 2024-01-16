package crawler

import (
	"fragrance-ws/models"
	"github.com/PuerkitoBio/goquery"
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
		return models.FragrancePage{}, err
	}

	fragrance := parseFragrance(doc)
	fragranceHouse := parseFragranceHouse(doc)
	noteCategories := parseNotePyramid(doc)
	accords := parseAccords(doc)

	return models.FragrancePage{
		Fragrance:   fragrance,
		FragHouse:   fragranceHouse,
		NoteCat:     noteCategories,
		MainAccords: accords,
	}, nil
}

func parseFragrance(doc *goquery.Document) models.Fragrance {
	name := doc.Find("h1").Text()
	description := doc.Find("div[itemprop='description'").Text()
	return models.Fragrance{
		Name:        name,
		Description: description,
	}
}

func parseFragranceHouse(doc *goquery.Document) models.FragranceHouse {
	name := doc.Find("h1").Text()
	return models.FragranceHouse{
		Name: name,
	}
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
