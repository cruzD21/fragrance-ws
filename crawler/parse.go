package crawler

import (
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

func parseFragrancePage(res *http.Response) (*Fragrance, error) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	frag := &Fragrance{
		Name:  doc.Find("h1").Text(),
		Brand: doc.Find("span[itemprop='name'][class='vote-button-name']").Text(),
	}

	doc.Find("h4 > b").Each(func(i int, s *goquery.Selection) {
		pyramidLevel := s.Text() //pyramid accords
		s.Parent().Next().Contents().Children().Each(func(i int, s2 *goquery.Selection) {
			note := s2.Children().Last().Last().Text()
			switch pyramidLevel {
			case "Top Notes":
				frag.TopNotes = append(frag.TopNotes, note)
			case "Middle Notes":
				frag.MiddleNotes = append(frag.MiddleNotes, note)
			case "Base Notes":
				frag.BaseNotes = append(frag.BaseNotes, note)
			}

		})

	})
	doc.Find(".accord-bar").Each(func(i int, s *goquery.Selection) {
		accord := s.Text()
		frag.MainAccords = append(frag.MainAccords, accord)
	})
	return frag, nil
}
