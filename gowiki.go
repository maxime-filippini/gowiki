package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type verb struct {
	infinitive   string
	aspect       string
	translations []string
}

func replaceSpecialCharacters(infinitive string) string {

	charMap := map[string]string{
		"ò": "o",
		"ì": "i",
		"í": "i",
		"ȅ": "e",
		"è": "e",
		"á": "a",
		"à": "a",
		"ȁ": "a",
	}

	for k, v := range charMap {
		infinitive = strings.ReplaceAll(infinitive, k, v)
	}

	return infinitive

}

func getVerbItems(verbItems *[]verb) colly.HTMLCallback {
	return func(e *colly.HTMLElement) {

		verbTitle := e.DOM.Parent().NextFilteredUntil("h3", "h2").Find("span[id=Verb]").Parent()
		verbInfo := verbTitle.NextFiltered("p")
		verbInfinitive := replaceSpecialCharacters(verbInfo.Find("strong.Latn.headword").Text())
		verbAspect := verbInfo.Find(".gender").Text()

		// Get the translations
		verbTranslations := make([]string, 0)

		verbTitle.Next().Next().Find("ol li").Each(func(_ int, s *goquery.Selection) {
			verbTranslations = append(verbTranslations, s.Text())
		})

		*verbItems = append(*verbItems, verb{infinitive: verbInfinitive, aspect: verbAspect, translations: verbTranslations})

	}
}

func main() {

	allVerbs := make([]verb, 0)
	c := colly.NewCollector(colly.AllowedDomains("en.wiktionary.org"))
	verbItemCollector := c.Clone()

	c.OnHTML(".mw-category-group ul li a[href]", func(e *colly.HTMLElement) {
		fmt.Printf("Visiting %s\n", e.Text)
		pageURL := e.Request.AbsoluteURL(e.Attr("href"))
		verbItemCollector.Visit(pageURL)
	})

	verbItemCollector.OnHTML("h2 .mw-headline[id=Serbo-Croatian]", getVerbItems(&allVerbs))

	c.Visit("https://en.wiktionary.org/wiki/Category:Serbo-Croatian_verbs")

	for _, element := range allVerbs {
		fmt.Println(element)
	}
}
