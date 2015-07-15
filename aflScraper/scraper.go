package aflScraper

import (
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brentonmcs/afl/aflDataAccess"
	"github.com/brentonmcs/afl/aflShared"
)

//ScrapePages scrapes the forecasts and the current prices
func ScrapePages() string {

	activeRounds := getactiveRounds()

	baseURI := "http://footyforecaster.com/AFL/RoundForecast/%d_Round_%d"
	if len(activeRounds) == 0 {
		activeRound := aflDataAccess.GetCurrentRound()
		if activeRound.Year != 0 {
			scrapeResults(fmt.Sprintf(baseURI, activeRound.Year, activeRound.Round+1), activeRound.Round+1, activeRound.Year)
		}
	}

	log.Print(activeRounds)

	scapeSportsBet()

	for _, aR := range activeRounds {

		fmt.Printf("scraping %d \n", aR.Round)
		scrapeResults(fmt.Sprintf(baseURI, aR.Year, aR.Round), aR.Round, aR.Year)
	}
	return "Done"
}

//SeedPages scrapes All forecasts and the current prices
func SeedPages() string {

	baseURI := "http://footyforecaster.com/AFL/RoundForecast/%d_Round_%d"

	for year := 2008; year < 2016; year++ {
		for round := 1; year < 24; round++ {
			scrapeResults(fmt.Sprintf(baseURI, year, round), round, year)
			fmt.Printf("Finished scraping year: %d round :%d", year, round)
		}
	}
	return "Done"
}

func goToPage(uri string) *goquery.Document {
	doc, err := goquery.NewDocument(uri)
	aflShared.HandleError(err)
	return doc
}

func scrapeResults(uri string, round, year int) {

	doc := goToPage(uri)

	updatedTime := time.Now()
	doc.Find(".details").Each(func(i int, s *goquery.Selection) {

		tableBody := s.Find("tbody")
		forecast := tableBody.Find("tr:nth-child(3) td:nth-child(2) td:first-child").Text()
		percentage := tableBody.Find("tr:nth-child(1) > td:nth-child(2)").Text()
		result := tableBody.Find("tr:nth-child(4) > td:nth-child(2) ").Text()

		var forecastModel = parseForecast(forecast, percentage, result, round, year, i, updatedTime)
		addForecast(&forecastModel)
	})

	removeOldForecasts(updatedTime, round, year)

}
